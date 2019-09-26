package core

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/mr-tron/base58/base58"
	shared "github.com/textileio/go-textile-core/bots"
	"github.com/textileio/go-textile/crypto"
)

type BotIpfsHandler struct {
	botID string
	node  *Textile
}

type BotKVStore struct {
	botID      string
	botVersion int
	store      map[string]string
	node       *Textile
}

func (mip BotIpfsHandler) Get(pth string, key string) ([]byte, error) {
	data, err := mip.node.DataAtPath(pth)

	if err != nil {
		return nil, err
	}

	// attempt decrypt if key present
	if key != "" {
		keyb, err := base58.Decode(key)
		if err != nil {
			log.Debugf("error decoding key %s: %s", key, err)
			return nil, err
		}
		plain, err := crypto.DecryptAES(data, keyb)
		if err != nil {
			log.Debugf("error decrypting %s: %s", pth, err)
			return nil, err
		}
		return plain, nil
	}
	return data, nil
}

func (kv BotKVStore) Set(key string, data []byte) (ok bool, err error) {
	err = kv.node.datastore.Bots().AddOrUpdate(kv.botID, key, data, kv.botVersion)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (kv BotKVStore) Get(key string) (data []byte, err error) {
	// TODO: include bot version from row in response, allowing migrations
	keyVal := kv.node.datastore.Bots().Get(kv.botID, key)
	if keyVal == nil {
		return []byte(""), nil
	}
	if keyVal.Value == nil {
		return []byte(""), nil
	}
	return keyVal.Value, nil
}
func (kv BotKVStore) Delete(key string) (ok bool, err error) {
	err = kv.node.datastore.Bots().Delete(kv.botID, key)
	if err != nil {
		return false, err
	}
	return true, nil
}

type BotService struct {
	clients map[string]*BotClient
	node    *Textile
}

// Exists is a helper to check if a bot exists
func (s *BotService) Exists(id string) bool {
	if s.clients == nil {
		return false
	}
	if _, ok := s.clients[id]; !ok {
		return false
	}
	return true
}

// Create configures the Bot rpc instance
func (s *BotService) Create(botID string, botVersion int, name string, pth string) {
	if s.Exists(botID) {
		return
	}

	store := &BotKVStore{
		botID,
		botVersion,
		make(map[string]string),
		s.node,
	}
	ipfs := &BotIpfsHandler{
		botID,
		s.node,
	}

	botClient := &BotClient{}
	s.clients[botID] = botClient
	s.clients[botID].setup(botID, botVersion, name, pth, store, ipfs)
}

// Get runs the bot.Get method
func (s *BotService) Get(botID string, q []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	botClient := s.clients[botID]
	res, err := botClient.service.Get(q, botClient.store, botClient.ipfs)
	return res, err
}

// Post runs the bot.Post method
func (s *BotService) Post(botID string, q []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	bot := s.clients[botID]
	res, err := bot.service.Post(q, bot.store, bot.ipfs)
	return res, err
}

// Put runs the bot.Put method
func (s *BotService) Put(botID string, q []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	bot := s.clients[botID]
	res, err := bot.service.Put(q, bot.store, bot.ipfs)
	return res, err
}

// Delete runs the bot.Delete method
func (s *BotService) Delete(botID string, q []byte) (shared.Response, error) {
	if !s.Exists(botID) {
		// TODO add error
		return shared.Response{
			Status: 400,
			Body:   []byte(""),
		}, nil
	}
	bot := s.clients[botID]
	res, err := bot.service.Delete(q, bot.store, bot.ipfs)
	return res, err
}

// RunAll runs a list of bots from Textile config
func (s *BotService) RunAll(repoPath string, bots []string) {
	for _, botID := range bots {
		botFolder := path.Join(repoPath, "bots")
		botPath := path.Join(botFolder, botID)
		botConfig, err := readBotConfig(botPath)
		if err != nil {
			log.Errorf(err.Error())
		} else {
			botPath := path.Join(botPath, "bot") // bots are always compiled to "bot"
			s.Create(botConfig.BotID, botConfig.ReleaseVersion, botConfig.BotName, botPath)
		}
	}
}

// NewBotService returns a new bot service
func NewBotService(node *Textile) *BotService {
	bots := &BotService{
		map[string]*BotClient{},
		node,
	}
	return bots
}

// ReadConfig loads the BotConfig
func readBotConfig(botPath string) (*shared.BotConfig, error) {
	data, err := ioutil.ReadFile(path.Join(botPath, "config"))
	if err != nil {
		return nil, err
	}

	var conf *shared.BotConfig
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}
