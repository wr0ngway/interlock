package dnsmasq

type PluginConfig struct {
	Port       int    `json:"port,omitempty"`
	PidPath    string `json:"pid_path,omitempty"`
	ConfigPath string `json:"config_path,omitempty`
	Domain     string `json:"domain,omitempty"`
}
