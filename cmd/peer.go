package cmd

func init() {
	register(&peerCmd{})
}

type peerCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *peerCmd) Name() string {
	return "peer"
}

func (x *peerCmd) Short() string {
	return "Show peer ID"
}

func (x *peerCmd) Long() string {
	return "Shows the local node's peer ID."
}

func (x *peerCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(GET, "peer", params{})
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
