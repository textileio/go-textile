package config

import (
	"encoding/json"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"io/ioutil"
	"os"
	"path"
)

var log = logging.Logger("tex-repo-config")

// Config is used to load textile config files.
type Config struct {
	Account   Account   // local node's account (public info only)
	Addresses Addresses // local node's addresses
	API       API       // local node's API settings
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
	API     string // address of the local API (RPC)
	CafeAPI string // address of the cafe REST API
	Gateway string // address to listen on for IPFS HTTP object gateway
}

type SwarmPorts struct {
	TCP string // TCP address port
	WS  string // WS address port
}

// API settings
type API struct {
	HTTPHeaders map[string][]string // HTTP headers to return with the API.
}

// Logs settings
type Logs struct {
	LogToDisk bool   // when true, sends all logs to rolling files on disk
	LogLevel  string // one of: critical error warning notice info debug
}

// Thread settings
type Threads struct {
	Default string // default thread to write to
}

// Cafe settings
// TODO: add some more knobs: max num. clients, max client msg age, inbox size, etc.
type Cafe struct {
	Open bool // when true, other peers can register with this node for cafe services
}

// Init returns the default textile config
func Init(version string) (*Config, error) {
	return &Config{
		Account: Account{
			Address: "",
			Thread:  "",
		},
		Addresses: Addresses{
			API:     "127.0.0.1:40600",
			CafeAPI: "127.0.0.1:40601",
			Gateway: "127.0.0.1:5050",
		},
		API: API{
			HTTPHeaders: map[string][]string{
				"Server": {"textile-go/" + version},
			},
		},
		Logs: Logs{
			LogToDisk: true,
			LogLevel:  "error",
		},
		Threads: Threads{
			Default: "",
		},
		Cafe: Cafe{
			Open: false,
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
