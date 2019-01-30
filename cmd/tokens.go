package cmd

import (
	"errors"
)

var errMissingToken = errors.New("missing token")

func init() {
	register(&tokensCmd{})
}

type tokensCmd struct {
	Create   createTokensCmd   `command:"create" description:"Create a new access token"`
	List     lsTokensCmd       `command:"ls" description:"List available access tokens"`
	Validate validateTokensCmd `command:"validate" description:"Check if access token is valid"`
	Remove   rmTokensCmd       `command:"rm" description:"Remove a specific access token"`
}

func (x *tokensCmd) Name() string {
	return "tokens"
}

func (x *tokensCmd) Short() string {
	return "Manage Cafe access tokens"
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

Generates an access token (44 random bytes) and saves a bcrypt hashed version for future lookup.
The response contains a base58 encoded version of the random bytes token.
`
}

func (x *createTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(POST, "tokens", params{})
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

List info about all stored cafe tokens.`
}

func (x *lsTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	var list []string
	res, err := executeJsonCmd(GET, "tokens", params{}, &list)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type validateTokensCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *validateTokensCmd) Usage() string {
	return `

Check validity of existing cafe access token.
`
}

func (x *validateTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) < 1 {
		return errMissingToken
	}
	res, err := executeStringCmd(GET, "tokens/"+args[0], params{})
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
	
	Removes an existing cafe token.`
}

func (x *rmTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) < 1 {
		return errMissingToken
	}
	res, err := executeStringCmd(DEL, "tokens/"+args[0], params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
