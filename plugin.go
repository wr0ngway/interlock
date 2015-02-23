package interlock

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Plugin interface {
		Name() string
		Version() string
		Url() string
		Author() string
		EventHandler(*Plugin, chan error)
		Info() *PluginInfo
	}

	PluginInput struct {
		Command string        `json:"command,omitempty"`
		Args    []interface{} `json:"args,omitempty"`
	}
)

func (p *Plugin) Start(ec chan error) (*File, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	// check if coming from pipe
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		go func() {
			for {
				err := <-ec
			}
		}()

		for {
			var input interlock.PluginInput
			if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
				return err
			}
			go p.EventHandler(&input, errorChan)
		}
	} else {
		// show version
	}
}
