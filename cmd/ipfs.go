package cmd

import (
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/textileio/textile-go/util"
)

var errMissingMultiAddress = errors.New("missing peer multi address")
var errMissingCID = errors.New("missing IPFS CID")

func init() {
	register(&ipfsCmd{})
}

type ipfsCmd struct {
	Swarm swarmCmd `command:"swarm" description:"Access some IPFS swarm commands"`
	Cat   catCmd   `command:"cat" description:"Show IPFS object data"`
}

func (x *ipfsCmd) Name() string {
	return "ipfs"
}

func (x *ipfsCmd) Short() string {
	return "Access IPFS commands"
}

func (x *ipfsCmd) Long() string {
	return "Provides access to some IPFS commands."
}

type swarmCmd struct {
	Connect swarmConnectCmd `command:"connect" description:"Open connection to a given address"`
	Peers   swarmPeersCmd   `command:"peers" description:"List peers with open connections"`
}

func (x *swarmCmd) Usage() string {
	return `

Opens a new direct connection to a peer address.
The address format is an IPFS multiaddr:

textile ipfs swarm connect /ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
`
}

type swarmConnectCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *swarmConnectCmd) Usage() string {
	return `

Opens a new direct connection to a peer address.

The address format is an IPFS multiaddr:

textile ipfs swarm connect /ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
`
}

func (x *swarmConnectCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingMultiAddress
	}

	res, err := executeJsonCmd(POST, "swarm/connect", params{
		args: args,
	}, nil)
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

	res, err := executeJsonCmd(GET, "swarm/peers", params{
		opts: map[string]string{
			"verbose":   strconv.FormatBool(x.Verbose),
			"streams":   strconv.FormatBool(x.Streams),
			"latency":   strconv.FormatBool(x.Latency),
			"direction": strconv.FormatBool(x.Direction),
		},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type catCmd struct {
	Client ClientOptions `group:"Client Options"`
	Key    string        `short:"k" long:"key" description:"Encyrption key."`
}

func (x *catCmd) Usage() string {
	return `

Displays the data behind an IPFS CID (hash).`
}

func (x *catCmd) Execute(args []string) error {
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
