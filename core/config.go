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
	conf, err := config.Read(init.RepoPath)
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
	return config.Write(init.RepoPath, conf)
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

// ensureMobileConfig ensures the low-power IPFS profile has been applied to the repo config
func ensureMobileConfig(repoPath string) error {
	rep, err := fsrepo.Open(repoPath)
	if err != nil {
		return err
	}
	conf, err := rep.Config()
	if err != nil {
		return err
	}

	conf.Routing.Type = "dhtclient"
	conf.Reprovider.Interval = "0"
	conf.Swarm.ConnMgr.LowWater = 200
	conf.Swarm.ConnMgr.HighWater = 500
	conf.Swarm.ConnMgr.GracePeriod = (time.Second * 20).String()
	conf.Swarm.DisableBandwidthMetrics = true
	conf.Swarm.EnableAutoRelay = true

	return rep.SetConfig(conf)
}

// ensureServerConfig ensures the server IPFS profile has been applied to the repo config
func ensureServerConfig(repoPath string) error {
	rep, err := fsrepo.Open(repoPath)
	if err != nil {
		return err
	}
	conf, err := rep.Config()
	if err != nil {
		return err
	}

	conf.Discovery.MDNS.Enabled = false
	conf.Addresses.NoAnnounce = config.DefaultServerFilters
	conf.Swarm.AddrFilters = config.DefaultServerFilters
	conf.Swarm.DisableNatPortMap = true
	conf.Swarm.EnableRelayHop = true
	conf.Swarm.EnableAutoNATService = true

	// tmp. ensure IPFS addresses are available in case we need to
	// point a vanilla daemon at the repo.
	conf.Addresses.API = []string{"/ip4/127.0.0.1/tcp/5001"}
	conf.Addresses.Gateway = []string{"/ip4/127.0.0.1/tcp/8080"}

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
