package interlock

import (
	"github.com/samalba/dockerclient"
)

type PluginAction struct {
	Name       string   `json:"name,omitempty"`
	EventName  string   `json:"event_name,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
}

type PluginInfo struct {
	Name        string          `json:"name,omitempty"`
	Version     string          `json:"version,omitempty"`
	Description string          `json:"description,omitempty"`
	Url         string          `json:"url,omitempty"`
	Actions     []*PluginAction `json:"actions,omitempty"`
}

type Plugin interface {
	Info() *PluginInfo
	Init() error
	HandleEvent(event *dockerclient.Event) error
}
