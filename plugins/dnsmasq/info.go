package dnsmasq

import (
	"github.com/ehazlett/interlock"
)

const (
	pluginName        = "dnsmasq"
	pluginVersion     = "0.1"
	pluginDescription = "dns services for containers"
	pluginUrl         = "https://github.com/ehazlett/interlock/tree/master/plugins/dnsmasq"
)

var (
	pluginInfo = &interlock.PluginInfo{
		Name:        pluginName,
		Version:     pluginVersion,
		Description: pluginDescription,
		Url:         pluginUrl,
	}
)
