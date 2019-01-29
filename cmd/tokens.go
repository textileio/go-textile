package cmd

import (
	"errors"

	"github.com/textileio/textile-go/repo"
)

var errMissingTokenId = errors.New("missing token id")
var errMissingToken = errors.New("missing token")

func init() {
	register(&tokensCmd{})
}

type tokensCmd struct {
	Create  createTokensCmd  `command:"create" description:"Create a new access token"`
	List    lsTokensCmd      `command:"ls" description:"List available access tokens"`
	Compare compareTokensCmd `command:"compare" description:"Check if access token is valid"`
	Remove  rmTokensCmd      `command:"rm" description:"Remove a specific access token"`
}

func (x *tokensCmd) Name() string {
	return "tokens"
}

func (x *tokensCmd) Short() string {
	return "Manage Cafe developer access tokens"
}

func (x *tokensCmd) Long() string {
	return `
Tokens allow other peers to register with a Cafe peer.
Use this command to create, list, compare, and remove tokens required for access to this peer's Cafe.
`
}

type createTokensCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *createTokensCmd) Usage() string {
	return `

Generates an access token (32 random bytes) and saves bcrypt encrypted version for future lookup.
The response contains a base58 encoded version of the random bytes token.
`
}

func (x *createTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	var info *repo.CafeDevToken
	res, err := executeJsonCmd(POST, "tokens", params{}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type lsTokensCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsTokensCmd) Usage() string {
	return `

List info about all stored cafe developer tokens.`
}

func (x *lsTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	var list []repo.CafeDevToken
	res, err := executeJsonCmd(GET, "tokens", params{}, &list)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type compareTokensCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *compareTokensCmd) Usage() string {
	return `

Check validity of existing cafe developer access token.
Requires a token id and the base58-encoded token itself.
`
}

func (x *compareTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) < 1 {
		return errMissingTokenId
	}
	if len(args) < 2 {
		return errMissingToken
	}
	res, err := executeStringCmd(GET, "tokens/"+args[0], params{
		args: []string{args[1]},
	})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type rmTokensCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmTokensCmd) Usage() string {
	return `
	
	Removes an existing cafe developer token.`
}

func (x *rmTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingTokenId
	}
	res, err := executeStringCmd(DEL, "tokens/"+args[0], params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
