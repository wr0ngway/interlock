package carbon

import (
	"github.com/ehazlett/interlock"
)

const (
	pluginName        = "carbon"
	pluginVersion     = "0.1"
	pluginDescription = "cluster stats to carbon (graphite)"
	pluginUrl         = "https://github.com/ehazlett/interlock/tree/master/plugins/carbon"
)

var (
	pluginInfo = &interlock.PluginInfo{
		Name:        pluginName,
		Version:     pluginVersion,
		Description: pluginDescription,
		Url:         pluginUrl,
	}
)
