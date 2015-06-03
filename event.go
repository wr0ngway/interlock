package interlock

import (
	"github.com/samalba/dockerclient"
)

type InterlockEvent struct {
	*dockerclient.Event
	Parameters map[string]string
}
