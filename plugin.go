package interlock

import (
	"github.com/samalba/dockerclient"
)

type PluginInfo struct {
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
}

type Plugin interface {
	Info() *PluginInfo
	Init() error
	HandleEvent(event *dockerclient.Event) error
}
