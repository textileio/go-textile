package core

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	bots "github.com/textileio/go-textile-bots"
	shared "github.com/textileio/go-textile-core/bots"
)

type BotClient struct {
	botID   string
	name    string
	service shared.Botservice
	config  *plugin.ClientConfig
	client  *plugin.Client
	store   *BotKVStore
	ipfs    *BotIpfsHandler
}

func (b *BotClient) setup(botID string, version int, name string, pth string, store *BotKVStore, ipfs *BotIpfsHandler) {
	pluginMap := map[string]plugin.Plugin{
		botID: &bots.TextileBot{}, // <- the TextileBot interface will always be the same.
	}

	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  uint(version),
		MagicCookieKey:   botID,
		MagicCookieValue: name,
	}
	// https://github.com/hashicorp/go-plugin/blob/master/client.go#L108
	// We're a host. Start by launching the plugin process.
	b.config = &plugin.ClientConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
		Cmd:             exec.Command("sh", "-c", pth),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Managed: true,
		Logger: hclog.New(&hclog.LoggerOptions{
			Output: hclog.DefaultOutput,
			Level:  hclog.Error,
			Name:   "plugin",
		}),
		// TODO add SecureConfig *SecureConfig for Hash-based bots
	}
	b.client = plugin.NewClient(b.config)
	b.botID = botID
	b.name = name
	b.store = store
	b.ipfs = ipfs
	// in go-plugin examples we need defer Kill, but because Managed: true, do we?
	// defer b.client.Kill()
	b.run()
}

func (b *BotClient) run() {
	// defer magicLink.client.Kill()
	// Connect via RPC
	rpcClient, err := b.client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	// Request the plugin
	raw, err := rpcClient.Dispense(b.botID)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	// We should have a Botservice store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	b.service = raw.(shared.Botservice)
}
