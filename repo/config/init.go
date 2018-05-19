package config

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/op/go-logging"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo"
	native "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/config"

	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	ci "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

var log = logging.MustGetLogger("config")

const (
	minPort = 1024
	maxPort = 49151
)

var textileBootstrapAddresses = []string{
	// cluster elastic ip, node 4
	"/ip4/35.169.206.101/tcp/4001/ipfs/QmP4S3UhmuEcGCHBGyG4zQVducj81YeNnvkCxUnZJrUopp",
	"/ip6/2600:1f18:6061:9403:8a36:fc7f:45be:2610/tcp/4001/ipfs/QmP4S3UhmuEcGCHBGyG4zQVducj81YeNnvkCxUnZJrUopp",

	// relay node 5
	"/ip4/34.201.54.67/tcp/4001/ipfs/QmTUvaGZqEu7qJw6DuTyhTgiZmZwdp7qN4FD4FFV3TGhjM",
	"/ip6/2600:1f18:6061:9403:b15e:b223:3c2e:1ee9/tcp/4001/ipfs/QmTUvaGZqEu7qJw6DuTyhTgiZmZwdp7qN4FD4FFV3TGhjM",
}

var RemoteRelayNode = "QmTUvaGZqEu7qJw6DuTyhTgiZmZwdp7qN4FD4FFV3TGhjM"

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

func Init(isMobile bool) (*native.Config, error) {
	identity, err := identityConfig()
	if err != nil {
		return nil, err
	}

	bootstrapPeers, err := native.DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}

	// add our own bootstrap peer
	for _, addr := range textileBootstrapAddresses {
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

	// some of the below are taken from the not-yet-released "lowpower" profile preset
	// TODO: profile with these setting on / off
	//if isMobile {
	//	reproviderInterval = "0"
	//	swarmConnMgrLowWater = 20
	//	swarmConnMgrHighWater = 40
	//	swarmConnMgrGracePeriod = time.Minute.String()
	//}

	conf := &native.Config{
		API: native.API{
			HTTPHeaders: map[string][]string{
				"Server": {"go-ipfs/" + native.CurrentVersionNumber},
			},
		},

		// setup the node's default addresses.
		// NOTE: two swarm listen addrs, one tcp, one utp.
		Addresses: addressesConfig(isMobile),

		Datastore: datastore,
		Bootstrap: native.BootstrapPeerStrings(bootstrapPeers),
		Identity:  identity,
		Discovery: native.Discovery{
			MDNS: native.MDNS{
				Enabled:  true,
				Interval: 10,
			},
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

func addressesConfig(isMobile bool) native.Addresses {
	swarmPort := getRandomPort()
	gatewayPort := getRandomPort()

	return native.Addresses{
		Swarm: []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", swarmPort),
			// "/ip4/0.0.0.0/udp/4002/utp", // disabled for now.
			fmt.Sprintf("/ip6/::/tcp/%d", swarmPort),
		},
		Announce:   []string{},
		NoAnnounce: []string{},
		Gateway:    fmt.Sprintf("127.0.0.1:%d", gatewayPort),
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

// identityConfig initializes a new identity.
func identityConfig() (native.Identity, error) {
	// TODO guard higher up
	ident := native.Identity{}

	log.Infof("generating Ed25519 keypair for peer identity...")
	sk, pk, err := ci.GenerateKeyPair(ci.Ed25519, 4096) // bits are ignored for ed25519, so use any
	if err != nil {
		return ident, err
	}

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := sk.Bytes()
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	log.Infof("new peer identity: %s\n", ident.PeerID)
	return ident, nil
}

func getRandomPort() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(maxPort-minPort) + minPort
}
