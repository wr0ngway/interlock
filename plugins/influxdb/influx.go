package influxdb

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	influx "github.com/influxdb/influxdb/client"
	"github.com/samalba/dockerclient"
)

const (
	defaultImageNameRegex = ".*"
)

var (
	errorChan = make(chan error)
)

type InfluxPlugin struct {
	interlockConfig *interlock.Config
	pluginConfig    *PluginConfig
	client          *dockerclient.DockerClient
	conn            *influx.Client
}

func init() {
	plugins.Register(
		pluginInfo.Name,
		&plugins.RegisteredPlugin{
			New: NewPlugin,
			Info: func() *interlock.PluginInfo {
				return pluginInfo
			},
		})
}

func loadPluginConfig() (*PluginConfig, error) {
	defaultImageNameFilter := regexp.MustCompile(defaultImageNameRegex)

	cfg := &PluginConfig{
		InfluxDbAddr:        "",
		InfluxDbName:        "stats",
		InfluxDbUser:        "",
		InfluxDbPass:        "",
		InfluxDbMeasurement: "stats",
		ImageNameFilter:     defaultImageNameFilter,
		Interval:            10,
	}

	// load custom config via environment
	influxAddr := os.Getenv("INFLUXDB_ADDR")
	if influxAddr != "" {
		cfg.InfluxDbAddr = influxAddr
	}

	influxName := os.Getenv("INFLUXDB_NAME")
	if influxName != "" {
		cfg.InfluxDbName = influxName
	}

	influxMeasurement := os.Getenv("INFLUXDB_MEASUREMENT")
	if influxMeasurement != "" {
		cfg.InfluxDbMeasurement = influxMeasurement
	}

	influxUser := os.Getenv("INFLUXDB_USER")
	if influxUser != "" {
		cfg.InfluxDbUser = influxUser
	}

	influxPass := os.Getenv("INFLUXDB_PASS")
	if influxPass != "" {
		cfg.InfluxDbPass = influxPass
	}

	imageNameFilter := os.Getenv("INFLUXDB_IMAGE_NAME_FILTER")
	if imageNameFilter != "" {
		// validate regex
		r, err := regexp.Compile(imageNameFilter)
		if err != nil {
			return nil, err
		}
		cfg.ImageNameFilter = r
	}

	interval := os.Getenv("INFLUXDB_INTERVAL")
	if interval != "" {
		i, err := strconv.Atoi(interval)
		if err != nil {
			return nil, err
		}
		cfg.Interval = i
	}

	return cfg, nil
}

func NewPlugin(interlockConfig *interlock.Config, client *dockerclient.DockerClient) (interlock.Plugin, error) {
	p := InfluxPlugin{interlockConfig: interlockConfig, client: client}
	cfg, err := loadPluginConfig()
	if err != nil {
		return nil, err
	}
	p.pluginConfig = cfg

	log.Debugf("connecting to influx: addr=%q user=%q pass=%q",
		p.pluginConfig.InfluxDbAddr,
		p.pluginConfig.InfluxDbUser,
		p.pluginConfig.InfluxDbPass,
	)

	host, err := url.Parse(p.pluginConfig.InfluxDbAddr)
	if err != nil {
		return nil, err
	}

	conf := influx.Config{
		URL:      *host,
		Username: p.pluginConfig.InfluxDbUser,
		Password: p.pluginConfig.InfluxDbPass,
	}

	conn, err := influx.NewClient(conf)
	if err != nil {
		return nil, err
	}

	p.conn = conn

	// handle errorChan
	go func() {
		for {
			err := <-errorChan
			plugins.Log(pluginInfo.Name,
				log.ErrorLevel,
				err.Error(),
			)
		}
	}()

	return p, nil
}

func (p InfluxPlugin) initialize() error {
	containers, err := p.client.ListContainers(false, false, "")
	if err != nil {
		return err
	}

	for _, c := range containers {
		if err := p.startStats(c.Id); err != nil {
			errorChan <- err
		}
	}

	plugins.Log(pluginInfo.Name, log.InfoLevel, fmt.Sprintf("sending stats every %d seconds", p.pluginConfig.Interval))

	return nil
}

func (p InfluxPlugin) Init() error {
	return nil
}

func (p InfluxPlugin) handleStats(id string, cb dockerclient.StatCallback, ec chan error, args ...interface{}) {
	go p.client.StartMonitorStats(id, cb, ec, args)
}

func (p InfluxPlugin) Info() *interlock.PluginInfo {
	return pluginInfo
}

func (p InfluxPlugin) sendStat(fields map[string]interface{}, tags map[string]string, t *time.Time) error {
	point := influx.Point{
		Measurement: p.pluginConfig.InfluxDbMeasurement,
		Fields:      fields,
		Tags:        tags,
		Time:        *t,
		Precision:   "s",
	}

	bp := influx.BatchPoints{
		Points: []influx.Point{
			point,
		},
		Database:        p.pluginConfig.InfluxDbName,
		RetentionPolicy: "default",
	}

	if _, err := p.conn.Write(bp); err != nil {
		return err
	}

	return nil
}

func (p InfluxPlugin) sendEventStats(id string, stats *dockerclient.Stats, ec chan error, args ...interface{}) {
	timestamp := time.Now()
	// report every n seconds
	rem := math.Mod(float64(timestamp.Second()), float64(p.pluginConfig.Interval))
	if rem != 0 {
		return
	}

	if len(id) > 12 {
		id = id[:12]
	}
	cInfo, err := p.client.InspectContainer(id)
	if err != nil {
		ec <- err
		return
	}

	cName := cInfo.Name[1:]

	log.Debug(cName)

	type containerStat struct {
		Key   string
		Value interface{}
		Tag   string
	}

	memPercent := float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit) * 100.0

	statData := []containerStat{
		{
			Key:   "totalUsage",
			Value: stats.CpuStats.CpuUsage.TotalUsage,
			Tag:   "cpu",
		},
		{
			Key:   "usageInKernelmode",
			Value: stats.CpuStats.CpuUsage.UsageInKernelmode,
			Tag:   "cpu",
		},
		{
			Key:   "usageInUsermode",
			Value: stats.CpuStats.CpuUsage.UsageInUsermode,
			Tag:   "cpu",
		},
		{
			Key:   "usage",
			Value: stats.MemoryStats.Usage,
			Tag:   "memory",
		},
		{
			Key:   "maxUsage",
			Value: stats.MemoryStats.MaxUsage,
			Tag:   "memory",
		},
		{
			Key:   "failcnt",
			Value: stats.MemoryStats.Failcnt,
			Tag:   "memory",
		},
		{
			Key:   "limit",
			Value: stats.MemoryStats.Limit,
			Tag:   "memory",
		},
		{
			Key:   "percent",
			Value: memPercent,
			Tag:   "memory",
		},
		{
			Key:   "rxBytes",
			Value: stats.NetworkStats.RxBytes,
			Tag:   "network",
		},
		{
			Key:   "rxPackets",
			Value: stats.NetworkStats.RxPackets,
			Tag:   "network",
		},
		{
			Key:   "rxErrors",
			Value: stats.NetworkStats.RxErrors,
			Tag:   "network",
		},
		{
			Key:   "rxDropped",
			Value: stats.NetworkStats.RxDropped,
			Tag:   "network",
		},
		{
			Key:   "txBytes",
			Value: stats.NetworkStats.TxBytes,
			Tag:   "network",
		},
		{
			Key:   "txPackets",
			Value: stats.NetworkStats.TxPackets,
			Tag:   "network",
		},
		{
			Key:   "txErrors",
			Value: stats.NetworkStats.TxErrors,
			Tag:   "network",
		},
		{
			Key:   "txDropped",
			Value: stats.NetworkStats.TxDropped,
			Tag:   "network",
		},
	}

	// send every n seconds
	for _, s := range statData {
		plugins.Log(pluginInfo.Name,
			log.DebugLevel,
			fmt.Sprintf("stat t=%d id=%s key=%s value=%v tags=%s",
				timestamp.UnixNano(),
				id,
				s.Key,
				s.Value,
				s.Tag,
			),
		)
		fields := map[string]interface{}{}
		fields[s.Key] = s.Value
		tags := map[string]string{
			"source": s.Tag,
		}
		if err := p.sendStat(fields, tags, &timestamp); err != nil {
			ec <- err
		}
	}

	return
}

func (p InfluxPlugin) startStats(id string) error {
	// get container info for event
	c, err := p.client.InspectContainer(id)
	if err != nil {
		return err
	}
	// match regex to start monitoring
	if p.pluginConfig.ImageNameFilter.MatchString(c.Config.Image) {
		plugins.Log(pluginInfo.Name, log.DebugLevel,
			fmt.Sprintf("gathering stats: image=%s id=%s", c.Image, c.Id[:12]))
		go p.handleStats(id, p.sendEventStats, errorChan, nil)
	}

	return nil
}

func (p InfluxPlugin) HandleEvent(event *interlock.InterlockEvent) error {
	// check all containers to see if stats are needed
	if err := p.initialize(); err != nil {
		return err
	}

	if event.Status == "start" {
		if err := p.startStats(event.Id); err != nil {
			return err
		}
	}
	return nil
}
