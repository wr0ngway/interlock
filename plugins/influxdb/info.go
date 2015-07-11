package influxdb

import (
	"github.com/ehazlett/interlock"
)

const (
	pluginName        = "influxdb"
	pluginVersion     = "0.1"
	pluginDescription = "cluster stats to influxdb"
	pluginUrl         = "https://github.com/ehazlett/interlock/tree/master/plugins/influxdb"
)

var (
	pluginInfo = &interlock.PluginInfo{
		Name:        pluginName,
		Version:     pluginVersion,
		Description: pluginDescription,
		Url:         pluginUrl,
	}
)
