package core

import (
	"fmt"
	"github.com/textileio/textile-go/repo/config"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo"
	"strings"
)

// Config returns the textile configuration file
func (t *Textile) Config() *config.Config {
	return t.config
}

// applyTextileConfigOptions update textile config w/ init options
func applyTextileConfigOptions(init InitConfig) error {
	conf, err := config.Read(init.RepoPath)
	if err != nil {
		return err
	}

	// account settings
	conf.Account.Address = init.Account.Address()

	// address settings
	conf.Addresses.API = init.ApiAddr
	conf.Addresses.CafeAPI = init.CafeApiAddr
	conf.Addresses.Gateway = init.GatewayAddr

	// log settings
	conf.Logs.LogToDisk = init.LogToDisk
	conf.Logs.LogLevel = strings.ToLower(init.LogLevel.String())

	// profile settings
	conf.IsServer = init.IsServer
	conf.IsMobile = init.IsMobile

	// cafe settings
	conf.Cafe.Open = init.CafeOpen

	// write to disk
	return config.Write(init.RepoPath, conf)
}

//func ensureBootstrapConfig(rep repo.Repo) error {
//	return config.UpdateIpfs(rep, "Bootstrap", config.BootstrapAddresses)
//	return nil
//}

// applySwarmPortConfigOption sets custom swarm ports (tcp and ws)
func applySwarmPortConfigOption(rep repo.Repo, ports string) error {
	parts := strings.Split(ports, ",")
	if len(parts) != 2 {
		return nil
	}
	return config.UpdateIpfs(rep, "Addresses.Swarm", []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", parts[0]),
		fmt.Sprintf("/ip6/::/tcp/%s", parts[0]),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", parts[1]),
		fmt.Sprintf("/ip6/::/tcp/%s/ws", parts[1]),
	})
	log.Infof("applied custom swarm port: %s", ports)
	return nil
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
