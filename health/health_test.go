package health_test

import (
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/mesher/cmd"
	"github.com/go-chassis/mesher/config"
	"github.com/go-chassis/mesher/health"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestHttpCheck(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	server := &http.Server{
		Addr: "127.0.0.1:3000",
	}
	http.HandleFunc("/health", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(200)
		resp.Write([]byte("hello"))
	})

	t.Log("Check server stoped")
	check := &config.HealthCheck{
		PortName: "rest-console",
		URI:      "/health",
	}
	addr := "127.0.0.1:3000"
	err := health.HTTPCheck(check, addr)
	assert.Error(t, err)

	t.Log("launch server")
	go server.ListenAndServe()
	time.Sleep(1000)
	defer server.Shutdown(nil)

	t.Log("check real health ")
	check = &config.HealthCheck{
		PortName: "rest-console",
		URI:      "/health",
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.NoError(t, err)

	t.Log("status match 500,must fail ")
	check = &config.HealthCheck{
		PortName: "rest-console",
		URI:      "/health",
		Match: &config.Match{
			Status: "201",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.Error(t, err)

	t.Log("body match fake,must fail ")
	check = &config.HealthCheck{
		PortName: "rest-console",
		URI:      "/health",
		Match: &config.Match{
			Body: "fake",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.Error(t, err)

	t.Log("body match right,no error ")
	check = &config.HealthCheck{
		PortName: "rest-console",
		URI:      "/health",
		Match: &config.Match{
			Body: "hello",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.NoError(t, err)

	t.Log("all match,no error ")
	check = &config.HealthCheck{
		PortName: "rest-console",
		URI:      "/health",
		Match: &config.Match{
			Status: "200",
			Body:   "hello",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.NoError(t, err)
}

func TestParseConfig(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	cmd.Configs = &cmd.ConfigFromCmd{
		PortsMap: map[string]string{
			"rest-console": "127.0.0.1:9000",
		},
	}
	check := &config.HealthCheck{
		PortName: "rest-console",
		URI:      "/health",
		Match: &config.Match{
			Body: "hello",
		},
	}
	addr, c, err := health.ParseConfig(check)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1:9000", addr)
	err = c(check, addr)
	assert.Error(t, err)
}
func TestL4Check(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	err := health.L4Check("127.0.0.1:3000")
	assert.Error(t, err)
	net.Listen("tcp", "127.0.0.1:3000")
	err = health.L4Check("127.0.0.1:3000")
	assert.NoError(t, err)
}
