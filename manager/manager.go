package manager

import (
	"crypto/tls"
	"net"
	"net/url"
	"os/exec"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	"github.com/samalba/dockerclient"
)

var (
	eventsErrChan = make(chan error)
	recovering    = false
)

type (
	Mgr struct {
		config    *interlock.Config
		client    *dockerclient.DockerClient
		mux       sync.Mutex
		tlsConfig *tls.Config
		proxyCmd  *exec.Cmd
	}
)

func NewManager(cfg *interlock.Config, tlsConfig *tls.Config) *Mgr {
	m := &Mgr{
		config:    cfg,
		tlsConfig: tlsConfig,
	}

	return m
}

func (m *Mgr) connect() error {
	log.Debugf("connecting to swarm on %s", m.config.SwarmUrl)

	c, err := dockerclient.NewDockerClient(m.config.SwarmUrl, m.tlsConfig)
	if err != nil {
		log.Warn(err)
		return err
	}

	m.client = c

	go m.startEventListener()

	return nil
}

func (m *Mgr) Config() *interlock.Config {
	return m.config
}

func (m *Mgr) Client() *dockerclient.DockerClient {
	return m.client
}

func (m *Mgr) startEventListener() {
	evt := NewEventHandler(m)
	m.client.StartMonitorEvents(evt.Handle, eventsErrChan)
}

func waitForTCP(addr string) error {
	log.Debugf("waiting for swarm to become available on %s", addr)
	for {
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err != nil {
			continue
		}
		conn.Close()
		break
	}
	return nil
}

func (m *Mgr) reconnectOnFail() {
	for {
		err := <-eventsErrChan
		log.Debugf("error on events chan: err=%s", err)

		if !recovering {
			recovering = true

			for {

				log.Warnf("error receiving events; attempting to reconnect")

				time.Sleep(1 * time.Second)

				u, err := url.Parse(m.config.SwarmUrl)
				if err != nil {
					log.Warnf("unable to parse Swarm URL: %s", err)
					continue
				}

				if u.Scheme != "unix" {
					if err := waitForTCP(u.Host); err != nil {
						log.Warnf("error connecting to Swarm: %s", err)
						continue
					}
				}

				if err := m.connect(); err == nil {
					log.Debugf("re-connected to Swarm: %s", m.config.SwarmUrl)
					recovering = false
					break
				}
			}
		}
	}
}

func (m *Mgr) Plugins() map[string]*plugins.RegisteredPlugin {
	// plugins
	allPlugins := plugins.GetPlugins()
	if len(allPlugins) == 0 || len(m.config.EnabledPlugins) == 0 {
		log.Warnf("no plugins enabled")
	}

	enabledPlugins := make(map[string]*plugins.RegisteredPlugin)

	for _, v := range m.config.EnabledPlugins {
		if p, ok := allPlugins[v]; ok {
			log.Infof("loading plugin name=%s version=%s",
				p.Info().Name,
				p.Info().Version)
			enabledPlugins[v] = p
		}
	}

	return enabledPlugins
}

func (m *Mgr) Run() error {
	if err := m.connect(); err != nil {
		return err
	}

	go m.reconnectOnFail()

	enabledPlugins := m.Plugins()
	plugins.SetEnabledPlugins(enabledPlugins)

	// custom event to signal startup
	m.Signal("interlock-start", nil)
	return nil
}

func (m *Mgr) Stop() error {
	// custom event to signal shutdown
	m.Signal("interlock-stop", nil)
	return nil
}

// Signal fires an event of the specified type
func (m *Mgr) Signal(action string, params map[string]string) error {
	evt := &interlock.InterlockEvent{
		Event: &dockerclient.Event{
			Id:     "",
			Status: action,
			From:   "interlock",
			Time:   time.Now().UnixNano(),
		},
		Parameters: params,
	}
	plugins.DispatchEvent(m.config, m.client, evt, eventsErrChan)
	return nil
}
