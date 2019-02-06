package cmd

import (
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/util"
)

var errMissingMultiAddress = errors.New("missing peer multi address")
var errMissingCID = errors.New("missing IPFS CID")

func init() {
	register(&swarmCmd{})
	register(&ipfsCatCmd{})
}

type swarmCmd struct {
	Connect swarmConnectCmd `command:"connect" description:"Open connection to a given address"`
	Peers   swarmPeersCmd   `command:"peers" description:"List peers with open connections"`
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

type swarmConnectCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *swarmConnectCmd) Usage() string {
	return `

Opens a new direct connection to a peer address.

The address format is an IPFS multiaddr:

ipfs swarm connect /ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
`
}

func (x *swarmConnectCmd) Execute(args []string) error {
	setApi(x.Client)
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
	output(res)
	return nil
}

type swarmPeersCmd struct {
	Client    ClientOptions `group:"Client Options"`
	Verbose   bool          `short:"v" long:"verbose" description:"Display all extra information."`
	Streams   bool          `short:"s" long:"streams" description:"Also list information about open streams for each peer."`
	Latency   bool          `short:"l" long:"latency" description:"Also list information about latency to each peer."`
	Direction bool          `short:"d" long:"direction" description:"Also list information about the direction of connection."`
}

func (x *swarmPeersCmd) Usage() string {
	return `

Lists the set of peers this node is connected to.`
}

func (x *swarmPeersCmd) Execute(args []string) error {
	setApi(x.Client)
	var info *ipfs.ConnInfos
	res, err := executeJsonCmd(GET, "swarm/peers", params{
		opts: map[string]string{
			"verbose":   strconv.FormatBool(x.Verbose),
			"streams":   strconv.FormatBool(x.Streams),
			"latency":   strconv.FormatBool(x.Latency),
			"direction": strconv.FormatBool(x.Direction),
		},
	}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type ipfsCatCmd struct {
	Client ClientOptions `group:"Client Options"`
	Key    string        `short:"k" long:"key" description:"Encyrption key."`
}

func (x *ipfsCatCmd) Name() string {
	return "cat"
}

func (x *ipfsCatCmd) Short() string {
	return "Show IPFS object data."
}

func (x *ipfsCatCmd) Long() string {
	return "Displays the data behind an IPFS CID (hash)."
}

func (x *ipfsCatCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingCID
	}

	res, _, err := request(GET, "ipfs/"+args[0], params{
		opts: map[string]string{"key": x.Key},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return err
		}
		return errors.New(body)
	}

	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		return err
	}

	return nil
}
