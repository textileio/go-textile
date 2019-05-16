package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/textileio/go-textile/common"
)

// Config is used to load textile config files.
type Config struct {
	Account   Account   // local node's account (public info only)
	Addresses Addresses // local node's addresses
	API       API       // local node's API settings
	Gateway   Gateway   // local node's Gateway settings
	Logs      Logs      // local node's log settings
	Threads   Threads   // local node's thread settings
	IsMobile  bool      // local node is setup for mobile
	IsServer  bool      // local node is setup for a server w/ a public IP
	Cafe      Cafe      // local node cafe settings
}

// Account store public account info
type Account struct {
	Address string // public key (seed is stored in the _possibly_ encrypted datastore)
	Thread  string // thread id of the default account thread used for sync between account peers
}

// Addresses stores the (string) bind addresses for the node.
type Addresses struct {
	API       string // bind address of the local REST API
	CafeAPI   string // bind address of the cafe REST API
	Gateway   string // bind address of the IPFS object gateway
	Profiling string // bind address of the profiling API
}

type SwarmPorts struct {
	TCP string // TCP address port
	WS  string // WS address port
}

// HTTPHeaders to customise things like COR
type HTTPHeaders = map[string][]string

// API settings
type API struct {
	HTTPHeaders HTTPHeaders
	SizeLimit   int64 // Maximum file size limit to accept for POST requests in bytes
}

// Gateway settings
type Gateway struct {
	HTTPHeaders HTTPHeaders
}

// Logs settings
type Logs struct {
	LogToDisk bool // when true, sends all logs to rolling files on disk
}

// Thread settings
type Threads struct {
	Defaults ThreadDefaults // default settings
}

// ThreadDefaults settings
type ThreadDefaults struct {
	ID string // default thread ID for reads/writes
}

// Cafe settings
type Cafe struct {
	Host   CafeHost
	Client CafeClient
}

// TODO: add some more knobs: max num. clients, max client msg age, inbox size, etc.
type CafeHost struct {
	Open        bool   // When true, other peers can register with this node for cafe services.
	URL         string // Override the resolved URL of this cafe, useful for load HTTPS and/or load balancers
	NeighborURL string // Specifies the URL of a secondary cafe. Must return cafe info.
	SizeLimit   int64  // Maximum file size limit to accept for POST requests in bytes.
}

// CafeClient settings
type CafeClient struct {
	Mobile MobileCafeClient
}

// MobileCafeClient settings
type MobileCafeClient struct {
	// messages w/ size less than limit will be handled by the p2p cafe service,
	// messages w/ size greater than limit will be handled by the mobile OS's background
	// upload service and the cafe HTTP API
	P2PWireLimit int // deprecated
}

// Init returns the default textile config
func Init() (*Config, error) {
	return &Config{
		Account: Account{
			Address: "",
			Thread:  "",
		},
		Addresses: Addresses{
			API:       "127.0.0.1:40600",
			CafeAPI:   "0.0.0.0:40601",
			Gateway:   "127.0.0.1:5050",
			Profiling: "127.0.0.1:6060",
		},
		API: API{
			HTTPHeaders: HTTPHeaders{
				"Server": {"go-textile/" + common.Version},
				// Explicitly allow all methods
				"Access-Control-Allow-Methods": {
					http.MethodConnect,
					http.MethodDelete,
					http.MethodGet,
					http.MethodHead,
					http.MethodOptions,
					http.MethodPatch,
					http.MethodPost,
					http.MethodPut,
					http.MethodTrace,
				},
				"Access-Control-Allow-Headers": {
					// rs/cors default headers
					"Origin",
					"Accept",
					"Content-Type",
					"X-Requested-With",
					// reason why this is here is unknown
					"Method",
					// textile custom headers
					"X-Textile-Args",
					"X-Textile-Opts",
				},
				"Access-Control-Allow-Origin": {
					"http://localhost:*",
					"http://127.0.0.1:*",
				},
			},
			SizeLimit: 0,
		},
		Gateway: Gateway{
			HTTPHeaders: HTTPHeaders{
				// Explicitly allow all methods
				"Access-Control-Allow-Methods": {
					http.MethodConnect,
					http.MethodDelete,
					http.MethodGet,
					http.MethodHead,
					http.MethodOptions,
					http.MethodPatch,
					http.MethodPost,
					http.MethodPut,
					http.MethodTrace,
				},
				// Explicitly allow all headers
				"Access-Control-Allow-Headers": {
					"*",
				},
				// Explicitly allow all origins
				"Access-Control-Allow-Origin": {
					"*",
				},
			},
		},
		Logs: Logs{
			LogToDisk: true,
		},
		Threads: Threads{
			Defaults: ThreadDefaults{
				ID: "",
			},
		},
		Cafe: Cafe{
			Host: CafeHost{
				Open:        false,
				URL:         "",
				NeighborURL: "",
				SizeLimit:   0,
			},
			Client: CafeClient{
				Mobile: MobileCafeClient{
					P2PWireLimit: 0,
				},
			},
		},
		IsMobile: false,
		IsServer: false,
	}, nil
}

// Read reads config from disk
func Read(repoPath string) (*Config, error) {
	data, err := ioutil.ReadFile(path.Join(repoPath, "textile"))
	if err != nil {
		return nil, err
	}

	var conf *Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// Write replaces the on-disk version of config with the given one
func Write(repoPath string, conf *Config) error {
	f, err := os.Create(path.Join(repoPath, "textile"))
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
