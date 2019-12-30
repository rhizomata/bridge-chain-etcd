package model

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// Config ..
type Config struct {
	Cluster  string
	Name     string
	Hostname string
	Port     uint
	DataDir  string
	EtcdUrls []string
	// HeartbeatInterval Heartbeat Interval
	HeartbeatInterval uint

	// CheckHeartbeatInterval Heartbeat check Interval
	CheckHeartbeatInterval uint

	// AliveThreasholdSecond Heartbeat time Threashold
	AliveThreasholdSeconds uint
}

// ParseFlagConfig ..
func ParseFlagConfig() (config *Config) {
	clusterName := flag.String("cluster", "cluster1", "name of cluster")
	name := flag.String("name", "bridge1", "name of etcd server")
	host := flag.String("exposed-host", "0.0.0.0", "host name/IP")
	port := flag.Uint("port", 8080, "liesten port for daemon")
	dataDir := flag.String("data-dir", "chain-data", "local data directory")
	etcdUrls := flag.String("etcd-urls", "", "etcd-urls,...")
	heartbeatInterval := flag.Uint("heartbeat-interval", 2, "heartbeat interval(seconds)")
	checkHeartbeatInterval := flag.Uint("heartbeat-check-interval", 3, "heartbeat check interval(seconds)")
	aliveThreasholdSeconds := flag.Uint("alive-threashold", 7, "alive threashold seconds")

	flag.Parse()

	config = new(Config)
	config.Cluster = *clusterName
	config.Name = *name
	config.Hostname = *host
	config.Port = *port
	config.DataDir = *dataDir + "/" + *name

	if !strings.Contains(",", *etcdUrls) {
		config.EtcdUrls = []string{*etcdUrls}
	} else {
		config.EtcdUrls = strings.Split(",", *etcdUrls)
	}

	config.HeartbeatInterval = *heartbeatInterval * uint(time.Second)
	config.CheckHeartbeatInterval = *checkHeartbeatInterval * uint(time.Second)
	config.AliveThreasholdSeconds = *aliveThreasholdSeconds

	return config
}

// GetDaemonAddr ..
func (config Config) GetDaemonAddr() string {
	return config.Hostname + ":" + fmt.Sprint(config.Port)
}

// GetDaemonURL http://{GetDaemonAddr}
func (config Config) GetDaemonURL() string {
	return "http://" + config.GetDaemonAddr()
}
