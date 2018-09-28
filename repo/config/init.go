package config

import (
	"fmt"
	"github.com/op/go-logging"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo"
	native "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/config"
	"math/rand"
	"time"
)

var log = logging.MustGetLogger("config")

const (
	minPort = 1024
	maxPort = 49151
)

var BootstrapAddresses = []string{
	"/ip4/54.89.183.252/tcp/46539/ipfs/QmVVKdbGw5VnNLWFw4YjsJUtwU3xPFemzbAMeCJtAtr2YF",  // us-east-1 (1)
	"/ip4/34.203.191.62/tcp/21117/ipfs/QmPZhrJ47ym4be69HUp3FC6DMV3H7kDUQTuaaYP3okpKdq",  // us-east-1 (2)
	"/ip4/54.175.223.66/tcp/38180/ipfs/Qmen1NEX4FcVsPVdVS9HK3ZjdymePSWtLwhpggrnfr6i7F",  // us-east-1 (3)
	"/ip4/18.213.220.205/tcp/18343/ipfs/QmarZawkUWeFw1zxyNJBSKQvt3HX3AWVPxkktNaU34ojuK", // us-east-1 (4)

	"/ip4/52.59.250.251/tcp/45785/ipfs/QmZcuJi2ctkkjwTF2HBAN9tU8w1Gkj1r2ZUm7tfcBuXSsf", // eu-central-1 (1)
	"/ip4/35.159.51.5/tcp/8367/ipfs/QmUoxL8BJYKLPVMkyPHERA2uy8ns3uKmLC3Kc4WQRP28Vt",    // eu-central-1 (2)
	"/ip4/35.158.110.50/tcp/42603/ipfs/QmXvhsrgmwr9zHQWgpokQCpuuBCGFffo79PHJSadFX6TjD", // eu-central-1 (3)
}

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

func Init(identity native.Identity) (*native.Config, error) {
	var bootstrapPeers []native.BootstrapPeer
	for _, addr := range BootstrapAddresses {
		p, err := native.ParseBootstrapPeer(addr)
		bootstrapPeers = append(bootstrapPeers, p)
		if err != nil {
			return nil, err
		}
	}

	datastore := defaultDatastoreConfig()

	reproviderInterval := "12h"
	swarmConnMgrLowWater := DefaultConnMgrLowWater
	swarmConnMgrHighWater := DefaultConnMgrHighWater
	swarmConnMgrGracePeriod := DefaultConnMgrGracePeriod.String()

	conf := &native.Config{
		// setup the node's default addresses.
		// NOTE: two swarm listen addrs, one tcp, one utp.
		Addresses: addressesConfig(),

		Datastore: datastore,
		Bootstrap: native.BootstrapPeerStrings(bootstrapPeers),
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
		},
		Reprovider: native.Reprovider{
			Interval: reproviderInterval,
			Strategy: "all",
		},
		Swarm: native.SwarmConfig{
			ConnMgr: native.ConnMgr{
				LowWater:    swarmConnMgrLowWater,
				HighWater:   swarmConnMgrHighWater,
				GracePeriod: swarmConnMgrGracePeriod,
				Type:        "basic",
			},
		},
		Experimental: native.Experiments{
			FilestoreEnabled:     false,
			ShardingEnabled:      false,
			Libp2pStreamMounting: true,
		},
	}

	return conf, nil
}

func Update(rep repo.Repo, key string, value interface{}) error {
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
	swarmPort := GetRandomPort()
	swarmWSPort := GetRandomPort()

	return native.Addresses{
		Swarm: []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", swarmPort),
			fmt.Sprintf("/ip6/::/tcp/%d", swarmPort),
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/ws", swarmWSPort),
			fmt.Sprintf("/ip6/::/tcp/%d/ws", swarmWSPort),
			// "/ip4/0.0.0.0/udp/4002/utp", // disabled for now.
		},
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

func GetRandomPort() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(maxPort-minPort) + minPort
}
