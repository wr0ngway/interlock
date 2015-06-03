package example

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	"github.com/samalba/dockerclient"
)

type ExamplePlugin struct {
	interlockConfig *interlock.Config
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

func NewPlugin(interlockConfig *interlock.Config, client *dockerclient.DockerClient) (interlock.Plugin, error) {
	return ExamplePlugin{interlockConfig: interlockConfig, client: client}, nil
}

func (p ExamplePlugin) Info() *interlock.PluginInfo {
	return pluginInfo
}

func (p ExamplePlugin) HandleEvent(event *interlock.InterlockEvent) error {
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

func (p ExamplePlugin) Init() error {
	return nil
}
