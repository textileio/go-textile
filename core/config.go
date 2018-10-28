package core

import (
	"fmt"
	"github.com/textileio/textile-go/repo/config"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo"
	"strings"
)

func ensureBootstrapConfig(rep repo.Repo) error {
	//return config.Update(rep, "Bootstrap", config.BootstrapAddresses)
	return nil
}

func applySwarmPortConfigOption(rep repo.Repo, ports string) error {
	parts := strings.Split(ports, ",")
	if len(parts) != 2 {
		return nil
	}
	return config.Update(rep, "Addresses.Swarm", []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", parts[0]),
		fmt.Sprintf("/ip6/::/tcp/%s", parts[0]),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", parts[1]),
		fmt.Sprintf("/ip6/::/tcp/%s/ws", parts[1]),
	})
	log.Infof("applied custom swarm port: %s", ports)
	return nil
}

func applyServerConfigOption(rep repo.Repo, isServer bool) error {
	if isServer {
		if err := config.Update(rep, "Addresses.NoAnnounce", config.DefaultServerFilters); err != nil {
			return err
		}
		if err := config.Update(rep, "Swarm.AddrFilters", config.DefaultServerFilters); err != nil {
			return err
		}
		if err := config.Update(rep, "Swarm.EnableRelayHop", true); err != nil {
			return err
		}
		if err := config.Update(rep, "Discovery.MDNS.Enabled", false); err != nil {
			return err
		}
		log.Info("applied server profile")
	} else {
		if err := config.Update(rep, "Addresses.NoAnnounce", []string{}); err != nil {
			return err
		}
		if err := config.Update(rep, "Swarm.AddrFilters", []string{}); err != nil {
			return err
		}
		if err := config.Update(rep, "Swarm.EnableRelayHop", false); err != nil {
			return err
		}
		if err := config.Update(rep, "Discovery.MDNS.Enabled", true); err != nil {
			return err
		}
	}
	return nil
}
