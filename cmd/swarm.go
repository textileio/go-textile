package cmd

import (
	"errors"
	"strconv"

	"github.com/textileio/textile-go/ipfs"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingMultiAddress = errors.New("missing peer multi address")

func init() {
	register(&swarmCmd{})
}

type swarmCmd struct {
	Connect swarmConnectCmd `command:"connect"`
	Peers   swarmPeersCmd   `command:"peers"`
}

func (x *swarmCmd) Name() string {
	return "swarm"
}

func (x *swarmCmd) Short() string {
	return "Access IPFS swarm commands"
}

func (x *swarmCmd) Long() string {
	return "Provides access to some IPFS swarm commands."
}

func (x *swarmCmd) Shell() *ishell.Cmd {
	return nil
}

type swarmConnectCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *swarmConnectCmd) Name() string {
	return "connect"
}

func (x *swarmConnectCmd) Short() string {
	return "Open connection to a given address"
}

func (x *swarmConnectCmd) Long() string {
	return `
Opens a new direct connection to a peer address.

The address format is an IPFS multiaddr:

ipfs swarm connect /ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
`
}

func (x *swarmConnectCmd) Execute(args []string) error {
	setApi(x.Client)
	return callSwarmConnect(args)
}

func (x *swarmConnectCmd) Shell() *ishell.Cmd {
	return nil
}

func callSwarmConnect(args []string) error {
	if len(args) == 0 {
		return errMissingMultiAddress
	}

	var info []string
	res, err := executeJsonCmd(POST, "swarm/connect", params{
		args: args,
	}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type swarmPeersCmd struct {
	Client    ClientOptions `group:"Client Options"`
	Verbose   bool          `short:"v" long:"verbose" description:"Display all extra information."`
	Streams   bool          `short:"s" long:"streams" description:"Also list information about open streams for each peer."`
	Latency   bool          `short:"l" long:"latency" description:"Also list information about latency to each peer."`
	Direction bool          `short:"d" long:"direction" description:"Also list information about the direction of connection."`
}

func (x *swarmPeersCmd) Name() string {
	return "peers"
}

func (x *swarmPeersCmd) Short() string {
	return "List peers with open connections"
}

func (x *swarmPeersCmd) Long() string {
	return "Lists the set of peers this node is connected to."
}

func (x *swarmPeersCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"verbose":   strconv.FormatBool(x.Verbose),
		"streams":   strconv.FormatBool(x.Streams),
		"latency":   strconv.FormatBool(x.Latency),
		"direction": strconv.FormatBool(x.Direction),
	}
	return callSwarmPeers(opts)
}

func (x *swarmPeersCmd) Shell() *ishell.Cmd {
	return nil
}

func callSwarmPeers(opts map[string]string) error {
	var info *ipfs.ConnInfos
	res, err := executeJsonCmd(GET, "swarm/peers", params{
		opts: opts,
	}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
