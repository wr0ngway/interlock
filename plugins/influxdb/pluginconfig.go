package influxdb

import (
	"regexp"
)

type PluginConfig struct {
	InfluxDbAddr        string
	InfluxDbName        string
	InfluxDbUser        string
	InfluxDbPass        string
	InfluxDbMeasurement string
	ImageNameFilter     *regexp.Regexp
	Interval            int
}
