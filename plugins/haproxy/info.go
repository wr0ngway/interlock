package haproxy

import (
	"github.com/ehazlett/interlock"
)

const (
	pluginName        = "haproxy"
	pluginVersion     = "0.1"
	pluginDescription = "haproxy load balancer and reverse proxy"
	pluginUrl         = "https://github.com/ehazlett/interlock/tree/master/plugins/haproxy"
)

var (
	pluginActions = []*interlock.PluginAction{
		{
			Name:       "reload",
			EventName:  "haproxy-reload",
			Parameters: nil,
		},
	}

	pluginInfo = &interlock.PluginInfo{
		Name:        pluginName,
		Version:     pluginVersion,
		Description: pluginDescription,
		Url:         pluginUrl,
		Actions:     pluginActions,
	}
)
