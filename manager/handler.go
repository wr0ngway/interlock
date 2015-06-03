package manager

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	"github.com/samalba/dockerclient"
)

type (
	EventHandler struct {
		Manager Manager
	}
)

func NewEventHandler(mgr Manager) *EventHandler {
	return &EventHandler{
		Manager: mgr,
	}
}

func (l *EventHandler) handle(e *interlock.InterlockEvent, ec chan error, args ...interface{}) {
	plugins.Log("interlock", log.DebugLevel,
		fmt.Sprintf("event: date=%d type=%s image=%s container=%s", e.Time, e.Status, e.From, e.Id))

	go plugins.DispatchEvent(l.Manager.Config(), l.Manager.Client(), e, ec)
}

func (l *EventHandler) Handle(e *dockerclient.Event, ec chan error, args ...interface{}) {
	evt := &interlock.InterlockEvent{
		e,
		nil,
	}

	l.handle(evt, ec, args)
}
