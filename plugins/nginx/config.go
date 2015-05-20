package nginx

type NginxConfig struct {
	PluginConfig `json:"plugin_config,omitempty"`
	Hosts        []*Host `json:"hosts,omitempty"`
}
