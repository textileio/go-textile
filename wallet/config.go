package wallet

import (
	"github.com/textileio/textile-go/repo/config"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo"
)

func ensureBootstrapConfig(rep repo.Repo) error {
	return config.Update(rep, "Bootstrap", config.BootstrapAddresses)
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
