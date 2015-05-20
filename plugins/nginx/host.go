package nginx

type Host struct {
	ServerNames        []string  `json:"server_names,omitempty"`
	Port               int       `json:"port,omitempty"`
	SSLPort            int       `json:"ssl_port,omitempty"`
	SSL                bool      `json:"ssl,omitempty"`
	SSLCert            string    `json:"ssl_cert,omitempty"`
	SSLCertKey         string    `json:"ssl_cert_key,omitempty"`
	SSLOnly            bool      `json:"ssl_only,omitempty"`
	Upstream           *Upstream `json:"upstream,omitempty"`
	WebsocketEndpoints []string  `json:"websocket_endpoints,omitempty"`
}
