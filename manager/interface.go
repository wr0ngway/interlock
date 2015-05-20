package manager

import (
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	"github.com/samalba/dockerclient"
)

type Manager interface {
	Run() error
	Stop() error
	Plugins() map[string]*plugins.RegisteredPlugin
	Config() *interlock.Config
	Client() *dockerclient.DockerClient
}
