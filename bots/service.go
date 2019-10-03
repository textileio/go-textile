package bots

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/mr-tron/base58/base58"
	tbots "github.com/textileio/go-textile-bots"
	shared "github.com/textileio/go-textile-core/bots"
	core "github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/crypto"
	ipfs "github.com/textileio/go-textile/ipfs"
)

// BotIpfsHandler implements shared.IpfsHandler. Extends it by hanging on the the botID
type BotIpfsHandler struct {
	botID string
	node  *core.Textile
}

// BotKVStore implements shared.BotStore. Extends it with BotID and BotVersion
type BotKVStore struct {
	botID      string
	botVersion int
	node       *core.Textile
}

// Get allows a bot to get IPFS data by the cid/path. Allows optional key for decryption on the fly
func (mip BotIpfsHandler) Get(pth string, key string) ([]byte, error) {
	data, err := mip.node.DataAtPath(pth)

	if err != nil {
		return nil, err
	}

	// attempt decrypt if key present
	if key != "" {
		keyb, err := base58.Decode(key)
		if err != nil {
			// log.Debugf("error decoding key %s: %s", key, err)
			return nil, err
		}
		plain, err := crypto.DecryptAES(data, keyb)
		if err != nil {
			// log.Debugf("error decrypting %s: %s", pth, err)
			return nil, err
		}
		return plain, nil
	}
	return data, nil
}

// Add allows a bot to add data to IPFS. currently it does not pin the data, only adds.
func (mip BotIpfsHandler) Add(data []byte, encrypt bool) (hash string, key string, err error) {
	var input []byte
	k := ""
	if encrypt {
		aes, err := crypto.GenerateAESKey()
		if err != nil {
			return "", "", err
		}
		input, err = crypto.EncryptAES(data, aes)
		if err != nil {
			return "", "", err
		}
		k = base58.FastBase58Encoding(aes)
	} else {
		input = data
	}
	r := bytes.NewReader(input)
	idp, err := ipfs.AddData(mip.node.Ipfs(), r, false, false)
	if err != nil {
		return "", "", err
	}
	return idp.Hash().B58String(), k, nil
}

// Set allows a bot to add a key-val to the store
func (kv BotKVStore) Set(key string, data []byte) (ok bool, err error) {
	datastore := kv.node.Datastore()
	err = datastore.Bots().AddOrUpdate(kv.botID, key, data, kv.botVersion)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Get allows a bot to get a value by string. It responds with the version of the bot that wrote the data.
func (kv BotKVStore) Get(key string) (data []byte, version int32, err error) {
	// TODO: include bot version from row in response, allowing migrations
	datastore := kv.node.Datastore()
	keyVal := datastore.Bots().Get(kv.botID, key)
	if keyVal == nil {
		return []byte(""), 0, nil
	}
	if keyVal.Value == nil {
		return []byte(""), 0, nil
	}
	return keyVal.Value, keyVal.BotReleaseVersion, nil
}

// Delete allows a bot to delete a value in the kv store
func (kv BotKVStore) Delete(key string) (ok bool, err error) {
	datastore := kv.node.Datastore()
	err = datastore.Bots().Delete(kv.botID, key)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Service holds a map to all running bots on this node
type Service struct {
	clients map[string]*tbots.Client
	node    *core.Textile
}

// Exists is a helper to check if a bot exists
func (s *Service) Exists(id string) bool {
	if s.clients == nil {
		return false
	}
	if _, ok := s.clients[id]; !ok {
		return false
	}
	return true
}

// Create configures the Bot rpc instance
func (s *Service) Create(botID string, botVersion int, name string, params map[string]string, pth string) {
	if s.Exists(botID) {
		return
	}

	store := &BotKVStore{
		botID,
		botVersion,
		s.node,
	}
	ipfs := &BotIpfsHandler{
		botID,
		s.node,
	}

	config := shared.ClientConfig{
		store,
		ipfs,
		params,
	}
	botClient := &tbots.Client{}
	s.clients[botID] = botClient
	s.clients[botID].Prepare(botID, botVersion, name, pth, config)
}

// Get runs the bot.Get method
func (s *Service) Get(botID string, q []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	botClient := s.clients[botID]
	res, err := botClient.Service.Get(q, botClient.SharedConf)
	return res, err
}

// Post runs the bot.Post method
func (s *Service) Post(botID string, q []byte, body []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	botClient := s.clients[botID]
	res, err := botClient.Service.Post(q, body, botClient.SharedConf)
	return res, err
}

// Put runs the bot.Put method
func (s *Service) Put(botID string, q []byte, body []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	botClient := s.clients[botID]
	res, err := botClient.Service.Put(q, body, botClient.SharedConf)
	return res, err
}

// Delete runs the bot.Delete method
func (s *Service) Delete(botID string, q []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		// TODO add error
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	botClient := s.clients[botID]
	res, err := botClient.Service.Delete(q, botClient.SharedConf)
	return res, err
}

// RunAll runs a list of bots from Textile config
func (s *Service) RunAll(repoPath string, bots []string) {
	for _, botID := range bots {
		botFolder := path.Join(repoPath, "bots")
		botPath := path.Join(botFolder, botID)
		botConfig, err := readBotConfig(botPath)
		if err != nil {
			// log.Errorf(err.Error("Bots: config read error"))
		} else {
			botPath := path.Join(botPath, "bot") // bots are always compiled to "bot"
			s.Create(botConfig.ID, botConfig.ReleaseVersion, botConfig.Name, botConfig.Params, botPath)
		}
	}
}

// NewService returns a new bot service
func NewService(node *core.Textile) *Service {
	bots := &Service{
		map[string]*tbots.Client{},
		node,
	}
	return bots
}

// ReadConfig loads the HostConfig
func readBotConfig(botPath string) (*shared.HostConfig, error) {
	data, err := ioutil.ReadFile(path.Join(botPath, "config"))
	if err != nil {
		return nil, err
	}

	var conf *shared.HostConfig
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}
