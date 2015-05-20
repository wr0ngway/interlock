package commands

import (
	"crypto/tls"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/api"
	"github.com/ehazlett/interlock/manager"
	"github.com/ehazlett/interlock/utils"
	"github.com/ehazlett/interlock/version"
)

var CmdStart = cli.Command{
	Name:   "start",
	Usage:  "Start Interlock",
	Action: cmdStart,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "api",
			Usage: "Enable API",
		},
		cli.StringFlag{
			Name:  "api-listen-addr",
			Usage: "API listen address",
			Value: ":8080",
		},
	},
}

func cmdStart(c *cli.Context) {
	swarmUrl := c.GlobalString("swarm-url")
	swarmTlsCaCert := c.GlobalString("swarm-tls-ca-cert")
	swarmTlsCert := c.GlobalString("swarm-tls-cert")
	swarmTlsKey := c.GlobalString("swarm-tls-key")
	allowInsecureTls := c.GlobalBool("swarm-allow-insecure")

	// only load env vars if no args
	// check environment for docker client config
	envDockerHost := os.Getenv("DOCKER_HOST")
	if swarmUrl == "" && envDockerHost != "" {
		swarmUrl = envDockerHost
	}

	// only load env vars if no args
	envDockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	envDockerTlsVerify := os.Getenv("DOCKER_TLS_VERIFY")
	if swarmTlsCaCert == "" && envDockerCertPath != "" && envDockerTlsVerify != "" {
		swarmTlsCaCert = filepath.Join(envDockerCertPath, "ca.pem")
		swarmTlsCert = filepath.Join(envDockerCertPath, "cert.pem")
		swarmTlsKey = filepath.Join(envDockerCertPath, "key.pem")
	}

	config := &interlock.Config{}
	config.SwarmUrl = swarmUrl
	config.EnabledPlugins = c.GlobalStringSlice("plugin")

	// load tlsconfig
	var tlsConfig *tls.Config
	if swarmTlsCaCert != "" && swarmTlsCert != "" && swarmTlsKey != "" {
		log.Infof("using tls for communication with swarm")
		caCert, err := ioutil.ReadFile(swarmTlsCaCert)
		if err != nil {
			log.Fatalf("error loading tls ca cert: %s", err)
		}

		cert, err := ioutil.ReadFile(swarmTlsCert)
		if err != nil {
			log.Fatalf("error loading tls cert: %s", err)
		}

		key, err := ioutil.ReadFile(swarmTlsKey)
		if err != nil {
			log.Fatalf("error loading tls key: %s", err)
		}

		cfg, err := utils.GetTLSConfig(caCert, cert, key, allowInsecureTls)
		if err != nil {
			log.Fatalf("error configuring tls: %s", err)
		}
		tlsConfig = cfg
	}

	m := manager.NewManager(config, tlsConfig)

	log.Infof("interlock running version=%s", version.FULL_VERSION)
	if err := m.Run(); err != nil {
		log.Fatal(err)
	}

	go func() {
		err := <-errChan
		log.Error(err)
	}()

	if c.Bool("api") {
		listenAddr := c.String("api-listen-addr")
		log.Debugf("enabling api: addr=%s", listenAddr)
		cfg := &api.ApiConfig{
			ListenAddr: listenAddr,
			Manager:    m,
		}

		a := api.NewApi(cfg)
		go func() {
			if err := a.Run(); err != nil {
				errChan <- err
			}
		}()
	}

	waitForInterrupt()

	log.Infof("shutting down")
	if err := m.Stop(); err != nil {
		log.Fatal(err)
	}
}
