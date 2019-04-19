package cmd

import "github.com/textileio/go-textile/pb"

func init() {
	register(&summaryCmd{})
}

type summaryCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *summaryCmd) Name() string {
	return "summary"
}

func (x *summaryCmd) Short() string {
	return "Get a summary of local data"
}

func (x *summaryCmd) Long() string {
	return "Get a summary of the local node's data."
}

func (x *summaryCmd) Execute(args []string) error {
	setApi(x.Client)
	var info pb.Summary
	res, err := executeJsonPbCmd(GET, "summary", params{}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
