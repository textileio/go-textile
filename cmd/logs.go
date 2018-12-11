package cmd

import "strconv"

func init() {
	register(&logsCmd{})
}

type logsCmd struct {
	Client    ClientOptions `group:"Client Options"`
	Subsystem string        `short:"s" long:"subsystem" description:"The subsystem logging identifier. Omit for all."`
	Level     string        `short:"l" long:"level" description:"The log level, with 'debug' the most verbose and 'critical' the least verbose."`
	All       bool          `short:"a" long:"all" description:"Whether to list/change only Textile subsystems (default: false), or all available subsystems (true)."`
}

func (x *logsCmd) Name() string {
	return "logs"
}

func (x *logsCmd) Short() string {
	return "List and control Textile subsystem logs."
}

func (x *logsCmd) Long() string {
	return `List or change the verbosity of one or all subsystems log output.

Use the --subsystem option to control a given subsystem's log level. Omit to list/change all subsystems.
Use the --level option to control the log level. One of: debug, info, warning, error, critical. Omit to get current level.
Use the --all flag to include all available subsystems. Omit to list/edit Textile subsystems only.

Textile logs piggyback on the IPFS event logs, so this command allows users to control specific subsystem logs for finer control.
`
}

func (x *logsCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"subsystem": x.Subsystem,
		"level":     x.Level,
		"all":       strconv.FormatBool(x.All),
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
	var info []map[string]string
	res, err := executeJsonCmd(method, "logs"+subsystem, params{opts: opts}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
