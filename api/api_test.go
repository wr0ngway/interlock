package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
)

type MockManager struct{}

func (m MockManager) Run() error {
	return nil
}

func (m MockManager) Stop() error {
	return nil
}

func (m MockManager) Plugins() map[string]*plugins.RegisteredPlugin {
	return nil
}

func (m MockManager) Config() *interlock.Config {
	return nil
}

func (m MockManager) Client() *dockerclient.DockerClient {
	return nil
}

func getTestApi() (*Api, error) {
	log.SetLevel(log.ErrorLevel)
	m := &MockManager{}

	cfg := &ApiConfig{
		ListenAddr: ":8080",
		Manager:    m,
	}

	return NewApi(cfg), nil
}

func TestApiGetIndex(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(api.apiIndex))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK, "expected response code 200")
}
