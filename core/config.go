package core

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/rs/cors"
	"github.com/textileio/go-textile/repo/config"
)

const minPort = 1024
const maxPort = 49151

// Config returns the textile configuration file
func (t *Textile) Config() *config.Config {
	return t.config
}

// GetRandomPort returns a port within the acceptable range
func GetRandomPort() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return strconv.Itoa(rand.Intn(maxPort-minPort) + minPort)
}

// applyTextileConfigOptions update textile config w/ init options
func applyTextileConfigOptions(init InitConfig) error {
	repo, err := init.Repo()
	if err != nil {
		return err
	}
	conf, err := config.Read(repo)
	if err != nil {
		return err
	}

	// determine the account thread id
	atid, err := init.Account.Id()
	if err != nil {
		return err
	}

	// account settings
	conf.Account.Address = init.Account.Address()
	conf.Account.Thread = atid.Pretty()

	// address settings
	if init.ApiAddr != "" {
		conf.Addresses.API = init.ApiAddr
	}
	if init.CafeApiAddr != "" {
		conf.Addresses.CafeAPI = init.CafeApiAddr
	}
	if init.GatewayAddr != "" {
		conf.Addresses.Gateway = init.GatewayAddr
	}
	if init.ProfilingAddr != "" {
		conf.Addresses.Profiling = init.ProfilingAddr
	}

	// log settings
	conf.Logs.LogToDisk = init.LogToDisk

	// profile settings
	conf.IsServer = init.IsServer
	conf.IsMobile = init.IsMobile

	// cafe settings
	conf.Cafe.Host.Open = init.CafeOpen
	conf.Cafe.Host.URL = init.CafeURL
	conf.Cafe.Host.NeighborURL = init.CafeNeighborURL

	// write to disk
	return config.Write(repo, conf)
}

// applySwarmPortConfigOption sets custom swarm ports (tcp and ws)
func applySwarmPortConfigOption(rep repo.Repo, ports string) error {
	var parts []string
	if ports != "" {
		parts = strings.Split(ports, ",")
	}
	var tcp, ws string

	switch len(parts) {
	case 1:
		tcp = parts[0]
	case 2:
		tcp = parts[0]
		ws = parts[1]
	default:
		tcp = GetRandomPort()
		ws = GetRandomPort()
	}

	list := []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", tcp),
		fmt.Sprintf("/ip6/::/tcp/%s", tcp),
	}
	if ws != "" {
		list = append(list, fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", ws))
		list = append(list, fmt.Sprintf("/ip6/::/tcp/%s/ws", ws))
	}

	return rep.SetConfigKey("Addresses.Swarm", list)
}

// profile defines config settings for different environments
type profile int

const (
	desktopProfile profile = iota
	mobileProfile
	serverProfile
)

// ensureProfile ensures the config settings are active for the selected profile
func ensureProfile(profile profile, repoPath string) error {
	rep, err := fsrepo.Open(repoPath)
	if err != nil {
		return err
	}
	conf, err := rep.Config()
	if err != nil {
		return err
	}

	if profile == serverProfile {
		conf.Addresses.NoAnnounce = config.DefaultServerFilters
		conf.Swarm.AddrFilters = config.DefaultServerFilters
	} else {
		conf.Addresses.NoAnnounce = make([]string, 0)
		conf.Swarm.AddrFilters = make([]string, 0)
	}
	// tmp. ensure IPFS addresses are available in case we need to
	// point a vanilla daemon at the repo.
	conf.Addresses.API = []string{"/ip4/127.0.0.1/tcp/5001"}
	conf.Addresses.Gateway = []string{"/ip4/127.0.0.1/tcp/8080"}

	if profile == serverProfile {
		conf.Discovery.MDNS.Enabled = false
	} else {
		conf.Discovery.MDNS.Enabled = true
	}

	if profile == serverProfile {
		conf.Routing.Type = "dht"
		conf.Reprovider.Interval = "12h"
	} else {
		conf.Routing.Type = "dhtclient"
		conf.Reprovider.Interval = "0"
	}

	if profile == mobileProfile {
		conf.Swarm.ConnMgr.LowWater = 200
		conf.Swarm.ConnMgr.HighWater = 500
		conf.Swarm.DisableBandwidthMetrics = true
	} else {
		conf.Swarm.ConnMgr.LowWater = 600
		conf.Swarm.ConnMgr.HighWater = 900
		conf.Swarm.DisableBandwidthMetrics = false
	}
	conf.Swarm.ConnMgr.GracePeriod = (time.Second * 20).String()

	if profile == serverProfile {
		conf.Swarm.DisableNatPortMap = true
		conf.Swarm.EnableRelayHop = true
		conf.Swarm.EnableAutoRelay = false
		conf.Swarm.EnableAutoNATService = true
	} else {
		conf.Swarm.DisableNatPortMap = false
		conf.Swarm.EnableRelayHop = false
		conf.Swarm.EnableAutoRelay = true
		conf.Swarm.EnableAutoNATService = false
	}
	conf.Swarm.DisableRelay = false

	return rep.SetConfig(conf)
}

// ConvertHeadersToCorsOptions converts http headers into the format that cors options accepts
func ConvertHeadersToCorsOptions(headers config.HTTPHeaders) cors.Options {
	options := cors.Options{}

	control, ok := headers["Access-Control-Allow-Origin"]
	if ok && len(control) > 0 {
		options.AllowedOrigins = control
	}

	control, ok = headers["Access-Control-Allow-Methods"]
	if ok && len(control) > 0 {
		options.AllowedMethods = control
	}

	control, ok = headers["Access-Control-Allow-Headers"]
	if ok && len(control) > 0 {
		options.AllowedHeaders = control
	}

	return options
}
