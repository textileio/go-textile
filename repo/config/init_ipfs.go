package config

import (
	"fmt"
	"time"

	native "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	ipfs "gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/repo"
)

// DefaultServerFilters has a list of non-routable IPv4 prefixes
// according to http://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
var DefaultServerFilters = []string{
	"/ip4/10.0.0.0/ipcidr/8",
	"/ip4/100.64.0.0/ipcidr/10",
	"/ip4/169.254.0.0/ipcidr/16",
	"/ip4/172.16.0.0/ipcidr/12",
	"/ip4/192.0.0.0/ipcidr/24",
	"/ip4/192.0.0.0/ipcidr/29",
	"/ip4/192.0.0.8/ipcidr/32",
	"/ip4/192.0.0.170/ipcidr/32",
	"/ip4/192.0.0.171/ipcidr/32",
	"/ip4/192.0.2.0/ipcidr/24",
	"/ip4/192.168.0.0/ipcidr/16",
	"/ip4/198.18.0.0/ipcidr/15",
	"/ip4/198.51.100.0/ipcidr/24",
	"/ip4/203.0.113.0/ipcidr/24",
	"/ip4/240.0.0.0/ipcidr/4",
}

// TextileBootstrapAddresses are the addresses of cafe nodes run by the Textile team.
var TextileBootstrapAddresses = []string{
	"/ip4/18.144.12.135/tcp/4001/ipfs/12D3KooWGBW3LfzypK3zgV4QxdPyUm3aEuwBDMKRRpCPm9FrJvar",  // us-west-1a
	"/ip4/13.57.23.210/tcp/4001/ipfs/12D3KooWQue2dSRqnZTVvikoxorZQ5Qyyug3hV65rYnWYpYsNMRE",   // us-west-1c
	"/ip4/18.221.167.133/tcp/4001/ipfs/12D3KooWERmHT6g4YkrPBTmhfDLjfi8b662vFCfvBXqzcdkPGQn1", // us-east-2a
	"/ip4/18.224.173.65/tcp/4001/ipfs/12D3KooWLh9Gd4C3knv4XqCyCuaNddfEoSLXgekVJzRyC5vsjv5d",  // us-east-2b
	"/ip4/35.180.16.103/tcp/4001/ipfs/12D3KooWDhSfXZCBVAK6SNQu7h6mfGCBJtjMS44PW5YA5YCjVmjB",  // eu-west-3a
	"/ip4/35.180.35.45/tcp/4001/ipfs/12D3KooWBCZEDkZ2VxdNYKLLUACWbXMvW9SpVbbvoFR9CtH4qJv9",   // eu-west-3b
	"/ip4/13.250.53.27/tcp/4001/ipfs/12D3KooWQ5MR9Ugz9HkVU3fYFbiWbQR4jxKJB66JoSY7nP5ShsqQ",   // ap-southeast-1a
	"/ip4/3.1.49.130/tcp/4001/ipfs/12D3KooWDWJ473M3fXMEcajbaGtqgr6i6SvDdh5Ru9i5ZzoJ9Qy8",     // ap-southeast-1b
}

// TextileBootstrapPeers returns the (parsed) set of Textile bootstrap peers.
func TextileBootstrapPeers() ([]native.BootstrapPeer, error) {
	ps, err := native.ParseBootstrapPeers(TextileBootstrapAddresses)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse hardcoded bootstrap peers: %s
This is a problem with the Textile codebase. Please report it to the dev team.`, err)
	}
	return ps, nil
}

// InitIpfs create the IPFS config file
func InitIpfs(identity native.Identity) (*native.Config, error) {
	ipfsPeers, err := native.DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}
	textilePeers, err := TextileBootstrapPeers()
	if err != nil {
		return nil, err
	}
	peers := append(textilePeers, ipfsPeers...)

	conf := &native.Config{
		API: native.API{
			HTTPHeaders: map[string][]string{
				"Server": {"go-ipfs/" + ipfs.CurrentVersionNumber},
			},
		},

		// setup the node's default addresses.
		// NOTE: two swarm listen addrs, one tcp, one utp.
		Addresses: addressesConfig(),

		Datastore: defaultDatastoreConfig(),
		Bootstrap: native.BootstrapPeerStrings(peers),
		Identity:  identity,
		Discovery: native.Discovery{
			MDNS: native.MDNS{
				Enabled:  true,
				Interval: 10,
			},
		},

		Routing: native.Routing{
			Type: "dht",
		},

		// setup the node mount points.
		Mounts: native.Mounts{
			IPFS: "/ipfs",
			IPNS: "/ipns",
		},

		Ipns: native.Ipns{
			ResolveCacheSize: 128,
		},

		Gateway: native.Gateway{
			RootRedirect: "",
			Writable:     false,
			PathPrefixes: []string{},
			HTTPHeaders: map[string][]string{
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"GET"},
				"Access-Control-Allow-Headers": {"X-Requested-With", "Range"},
			},
			APICommands: []string{},
		},
		Reprovider: native.Reprovider{
			Interval: "12h",
			Strategy: "all",
		},
		Swarm: native.SwarmConfig{
			ConnMgr: native.ConnMgr{
				LowWater:    DefaultConnMgrLowWater,
				HighWater:   DefaultConnMgrHighWater,
				GracePeriod: DefaultConnMgrGracePeriod.String(),
				Type:        "basic",
			},
		},
		Experimental: native.Experiments{
			FilestoreEnabled:     false,
			ShardingEnabled:      false,
			Libp2pStreamMounting: true,
		},
		Pubsub: native.PubsubConfig{
			Router: "gossipsub",
		},
	}

	return conf, nil
}

// UpdateIpfs updates the IPFS config file
func UpdateIpfs(rep repo.Repo, key string, value interface{}) error {
	if err := rep.SetConfigKey(key, value); err != nil {
		log.Errorf("error setting %s: %s", key, err)
		return err
	}
	return nil
}

// DefaultConnMgrHighWater is the default value for the connection managers
// 'high water' mark
const DefaultConnMgrHighWater = 900

// DefaultConnMgrLowWater is the default value for the connection managers 'low
// water' mark
const DefaultConnMgrLowWater = 600

// DefaultConnMgrGracePeriod is the default value for the connection managers
// grace period
const DefaultConnMgrGracePeriod = time.Second * 20

func addressesConfig() native.Addresses {
	return native.Addresses{
		Swarm:      []string{},
		Announce:   []string{},
		NoAnnounce: []string{},
	}
}

// DefaultDatastoreConfig is an internal function exported to aid in testing.
func defaultDatastoreConfig() native.Datastore {
	return native.Datastore{
		StorageMax:         "10GB",
		StorageGCWatermark: 90, // 90%
		GCPeriod:           "1h",
		BloomFilterSize:    0,
		Spec: map[string]interface{}{
			"type": "mount",
			"mounts": []interface{}{
				map[string]interface{}{
					"mountpoint": "/blocks",
					"type":       "measure",
					"prefix":     "flatfs.datastore",
					"child": map[string]interface{}{
						"type":      "flatfs",
						"path":      "blocks",
						"sync":      true,
						"shardFunc": "/repo/flatfs/shard/v1/next-to-last/2",
					},
				},
				map[string]interface{}{
					"mountpoint": "/",
					"type":       "measure",
					"prefix":     "leveldb.datastore",
					"child": map[string]interface{}{
						"type":        "levelds",
						"path":        "datastore",
						"compression": "none",
					},
				},
			},
		},
	}
}
