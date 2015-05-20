package main

// interlock plugins
import (
	_ "github.com/ehazlett/interlock/plugins/carbon"
	_ "github.com/ehazlett/interlock/plugins/example"
	_ "github.com/ehazlett/interlock/plugins/haproxy"
	_ "github.com/ehazlett/interlock/plugins/nginx"
)
