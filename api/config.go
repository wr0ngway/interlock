package api

import (
	"github.com/ehazlett/interlock/manager"
)

type ApiConfig struct {
	ListenAddr string
	Manager    manager.Manager
}
