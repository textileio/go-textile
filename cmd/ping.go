package cmd

func init() {
	register(&pingCmd{})
}

type pingCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *pingCmd) Name() string {
	return "ping"
}

func (x *pingCmd) Short() string {
	return "Ping another peer"
}

func (x *pingCmd) Long() string {
	return "Pings another peer on the network, returning online|offline."
}

func (x *pingCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(GET, "ping", params{args: args})
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
