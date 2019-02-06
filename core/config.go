package core

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/repo"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/repo/fsrepo"

	"github.com/textileio/textile-go/repo/config"
)

const minPort = 1024
const maxPort = 49151

var tcpPortRx = regexp.MustCompile("/tcp/([0-9]+)$")
var wsPortRx = regexp.MustCompile("/tcp/([0-9]+)/ws$")

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
	pk, err := init.Account.LibP2PPubKey()
	if err != nil {
		return err
	}
	atid, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return err
	}

	// account settings
	conf.Account.Address = init.Account.Address()
	conf.Account.Thread = atid.Pretty()

	// address settings
	conf.Addresses.API = init.ApiAddr
	conf.Addresses.CafeAPI = init.CafeApiAddr
	conf.Addresses.Gateway = init.GatewayAddr

	// log settings
	conf.Logs.LogToDisk = init.LogToDisk

	// profile settings
	conf.IsServer = init.IsServer
	conf.IsMobile = init.IsMobile

	// cafe settings
	conf.Cafe.Host.Open = init.CafeOpen
	conf.Cafe.Host.PublicIP = init.CafePublicIP
	conf.Cafe.Host.URL = init.CafeURL
	conf.Cafe.Host.NeighborURL = init.CafeNeighborURL

	// write to disk
	return config.Write(init.RepoPath, conf)
}

// updateBootstrapConfig adds additional peers to the bootstrap config
func updateBootstrapConfig(repoPath string, add []string, rm []string) error {
	rep, err := fsrepo.Open(repoPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := rep.Close(); err != nil {
			log.Error(err)
		}
	}()
	conf, err := rep.Config()
	if err != nil {
		return err
	}
	var final []string

	// get a list that does not include items in rm
outer:
	for _, bp := range conf.Bootstrap {
		for _, r := range rm {
			if bp == r {
				continue outer
			}
		}
		final = append(final, bp)
	}

	for _, p := range add {
		final = append(final, p)
	}
	return config.UpdateIpfs(rep, "Bootstrap", final)
}

// loadSwarmPorts returns the swarm ports in the ipfs config
func loadSwarmPorts(repoPath string) (*config.SwarmPorts, error) {
	rep, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rep.Close(); err != nil {
			log.Error(err)
		}
	}()

	conf, err := rep.Config()
	if err != nil {
		return nil, err
	}
	ports := &config.SwarmPorts{}

	for _, p := range conf.Addresses.Swarm {
		tcp := tcpPortRx.FindStringSubmatch(p)
		if len(tcp) == 2 {
			ports.TCP = tcp[1]
		}
		ws := wsPortRx.FindStringSubmatch(p)
		if len(ws) == 2 {
			ports.WS = ws[1]
		}
	}
	return ports, nil
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
	}

	list := []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", tcp),
		fmt.Sprintf("/ip6/::/tcp/%s", tcp),
	}
	if ws != "" {
		list = append(list, fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", ws))
		list = append(list, fmt.Sprintf("/ip6/::/tcp/%s/ws", ws))
	}

	return config.UpdateIpfs(rep, "Addresses.Swarm", list)
}

// applyServerConfigOption adds the IPFS server profile to the repo config
func applyServerConfigOption(rep repo.Repo, isServer bool) error {
	if isServer {
		if err := config.UpdateIpfs(rep, "Addresses.NoAnnounce", config.DefaultServerFilters); err != nil {
			return err
		}
		if err := config.UpdateIpfs(rep, "Swarm.AddrFilters", config.DefaultServerFilters); err != nil {
			return err
		}
		if err := config.UpdateIpfs(rep, "Swarm.EnableRelayHop", true); err != nil {
			return err
		}
		if err := config.UpdateIpfs(rep, "Discovery.MDNS.Enabled", false); err != nil {
			return err
		}
		log.Info("applied server profile")

	} else {
		if err := config.UpdateIpfs(rep, "Addresses.NoAnnounce", []string{}); err != nil {
			return err
		}
		if err := config.UpdateIpfs(rep, "Swarm.AddrFilters", []string{}); err != nil {
			return err
		}
		if err := config.UpdateIpfs(rep, "Swarm.EnableRelayHop", false); err != nil {
			return err
		}
		if err := config.UpdateIpfs(rep, "Discovery.MDNS.Enabled", true); err != nil {
			return err
		}
	}

	return nil
}
