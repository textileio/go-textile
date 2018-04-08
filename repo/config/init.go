package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	native "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/config"

	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	ci "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

var textileBootstrapAddresses = []string{
	// cluster elastic ip, node 4
	"/ip4/35.169.206.101/tcp/4001/ipfs/QmP4S3UhmuEcGCHBGyG4zQVducj81YeNnvkCxUnZJrUopp",
	"/ip6/2600:1f18:6061:9403:8a36:fc7f:45be:2610/tcp/4001/ipfs/QmP4S3UhmuEcGCHBGyG4zQVducj81YeNnvkCxUnZJrUopp",

	// relay node 5
	"/ip4/34.201.54.67/tcp/4001/ipfs/QmTUvaGZqEu7qJw6DuTyhTgiZmZwdp7qN4FD4FFV3TGhjM",
	"/ip6/2600:1f18:6061:9403:b15e:b223:3c2e:1ee9/tcp/4001/ipfs/QmTUvaGZqEu7qJw6DuTyhTgiZmZwdp7qN4FD4FFV3TGhjM",
}

func Init(out io.Writer, nBitsForKeypair int, isMobile bool) (*native.Config, error) {
	identity, err := identityConfig(out, nBitsForKeypair)
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
	if isMobile {
		reproviderInterval = "0"
		swarmConnMgrLowWater = 20
		swarmConnMgrHighWater = 40
		swarmConnMgrGracePeriod = time.Minute.String()
	}

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
				"Access-Control-Allow-Origin":  []string{"*"},
				"Access-Control-Allow-Methods": []string{"GET"},
				"Access-Control-Allow-Headers": []string{"X-Requested-With", "Range"},
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
	swarmPort := "4002"
	apiPost := "5002"
	gatewayPort := "8081"
	if isMobile {
		swarmPort = "4003"
		apiPost = "5003"
		gatewayPort = "8082"
	}

	return native.Addresses{
		Swarm: []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", swarmPort),
			// "/ip4/0.0.0.0/udp/4002/utp", // disabled for now.
			fmt.Sprintf("/ip6/::/tcp/%s", swarmPort),
		},
		Announce:   []string{},
		NoAnnounce: []string{},
		API:        fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", apiPost),
		Gateway:    fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", gatewayPort),
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
func identityConfig(out io.Writer, nbits int) (native.Identity, error) {
	// TODO guard higher up
	ident := native.Identity{}
	if nbits < 1024 {
		return ident, errors.New("bitsize less than 1024 is considered unsafe")
	}

	fmt.Fprintf(out, "generating %v-bit RSA keypair...", nbits)
	sk, pk, err := ci.GenerateKeyPair(ci.RSA, nbits)
	if err != nil {
		return ident, err
	}
	fmt.Fprint(out, "done\n")

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
	fmt.Fprintf(out, "peer identity: %s\n", ident.PeerID)
	return ident, nil
}
