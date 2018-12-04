package cmd

func init() {
	register(&addressCmd{})
}

type addressCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *addressCmd) Name() string {
	return "address"
}
func (x *addressCmd) Short() string {
	return "Show wallet address"
}
func (x *addressCmd) Long() string {
	return "Shows the local node's wallet address."
}

func (x *addressCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(GET, "address", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
