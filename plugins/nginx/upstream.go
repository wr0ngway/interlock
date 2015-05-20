package nginx

type Server struct {
	Addr string `json:"addr,omitempty"`
}

type Upstream struct {
	Name    string    `json:"name,omitempty"`
	Servers []*Server `json:"servers,omitempty"`
}
