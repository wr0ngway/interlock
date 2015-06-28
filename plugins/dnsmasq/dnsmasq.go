package dnsmasq

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	"github.com/samalba/dockerclient"
)

type DNSMasqPlugin struct {
	interlockConfig *interlock.Config
	pluginConfig    *PluginConfig
	client          *dockerclient.DockerClient
}

func init() {
	plugins.Register(
		pluginInfo.Name,
		&plugins.RegisteredPlugin{
			New: NewPlugin,
			Info: func() *interlock.PluginInfo {
				return pluginInfo
			},
		})
}

func loadPluginConfig() (*PluginConfig, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg := &PluginConfig{
		Port:    80,
		PidPath: filepath.Join(wd, "dnsmasq.pid"),
		Domain:  "local",
	}

	// load custom config via environment
	configPath := os.Getenv("DNSMASQ_CONFIG_PATH")
	if configPath != "" {
		cfg.ConfigPath = configPath
	}

	domain := os.Getenv("DNSMASQ_DOMAIN")
	if domain != "" {
		cfg.Domain = domain
	}

	port := os.Getenv("DNSMASQ_PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		cfg.Port = p
	}

	pidPath := os.Getenv("DNSMASQ_PID_PATH")
	if pidPath != "" {
		cfg.PidPath = pidPath
	}

	return cfg, nil
}

func NewPlugin(interlockConfig *interlock.Config, client *dockerclient.DockerClient) (interlock.Plugin, error) {
	pluginConfig, err := loadPluginConfig()
	if err != nil {
		return nil, err
	}

	plugin := DNSMasqPlugin{
		pluginConfig:    pluginConfig,
		interlockConfig: interlockConfig,
		client:          client,
	}

	return plugin, nil
}

func (p DNSMasqPlugin) Info() *interlock.PluginInfo {
	return pluginInfo
}

func (p DNSMasqPlugin) HandleEvent(event *interlock.InterlockEvent) error {
	plugins.Log(pluginInfo.Name, log.InfoLevel,
		fmt.Sprintf("action=received cid=%q event=%s time=%d params=%q",
			event.Id,
			event.Status,
			event.Time,
			event.Parameters,
		),
	)
	return nil
}

func (p DNSMasqPlugin) Init() error {
	return nil
}
