package cmd

import "strconv"

func init() {
	register(&logsCmd{})
}

type logsCmd struct {
	Client    ClientOptions `group:"Client Options"`
	Subsystem string        `short:"s" long:"subsystem" description:"The subsystem logging identifier. Omit for all."`
	Level     string        `short:"l" long:"level" description:"One of: debug, info, warning, error, critical. Omit to get current level."`
	TexOnly   bool          `short:"t" long:"tex-only" description:"Whether to list/change only Textile subsystems, or all available subsystems."`
}

func (x *logsCmd) Name() string {
	return "logs"
}

func (x *logsCmd) Short() string {
	return "List and control Textile subsystem logs."
}

func (x *logsCmd) Long() string {
	return `
List or change the verbosity of one or all subsystems log output.
Textile logs piggyback on the IPFS event logs.
`
}

func (x *logsCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"subsystem": x.Subsystem,
		"level":     x.Level,
		"tex-only":  strconv.FormatBool(x.TexOnly),
	}
	return callLogs(opts)
}

func callLogs(opts map[string]string) error {
	subsystem := opts["subsystem"]
	if subsystem != "" {
		subsystem = "/" + subsystem
	}
	method := GET
	if opts["level"] != "" {
		method = POST
	}

	res, err := executeJsonCmd(method, "logs"+subsystem, params{opts: opts}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
