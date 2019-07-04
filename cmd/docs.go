package cmd

import (
	"fmt"
	"html"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Optional interface for flags that can be repeated.
type repeatableFlag interface {
	kingpin.Value
	IsCumulative() bool
}

func formatFlag(haveShort bool, flag *kingpin.FlagModel) string {
	flagString := ""
	if flag.Short != 0 {
		flagString += fmt.Sprintf("-%c, --%s", flag.Short, flag.Name)
	} else {
		if haveShort {
			flagString += fmt.Sprintf("    --%s", flag.Name)
		} else {
			flagString += fmt.Sprintf("--%s", flag.Name)
		}
	}
	if !flag.IsBoolFlag() {
		flagString += fmt.Sprintf("=%s", flag.FormatPlaceHolder())
	}
	if v, ok := flag.Value.(repeatableFlag); ok && v.IsCumulative() {
		flagString += " ..."
	}
	return flagString
}

func formatFlags(flags []*kingpin.FlagModel) string {
	var line = ""
	if len(flags) != 0 {
		line += "\n\t<tr><th>Flag</th><th>Description</th></tr>"
		for _, flag := range flags {
			name := formatFlag(flag.Short != 0, flag)
			desc := strings.TrimSpace(flag.Help)
			line += "\n\t<tr><td><code>" + html.EscapeString(name) + "</code></td><td><pre>" + html.EscapeString(desc) + "</pre></td></tr>"
		}
	}
	return line
}

func formatArgs(args []*kingpin.ArgModel) string {
	var line = ""
	if len(args) != 0 {
		line += "\n\t<tr><th>Argument</th><th>Description</th></tr>"
		for _, arg := range args {
			name := arg.Name
			if len(arg.Default) != 0 {
				name = name + "=" + strings.Join(arg.Default, ",")
			}
			name = "<" + name + ">"
			if !arg.Required {
				name = "[" + name + "]"
			}
			desc := strings.TrimSpace(arg.Help)
			line += "\n\t<tr><td><code>" + html.EscapeString(name) + "</code></td><td><pre>" + html.EscapeString(desc) + "</pre></td></tr>"
		}
	}
	return line
}

func formatCommands(cmds []*kingpin.CmdModel) string {
	var line = ""
	if len(cmds) != 0 {
		for _, cmd := range cmds {
			line += "\n" + formatCommand(appCmd, *cmd)
		}
	}
	return line
}

func formatCommand(appCmd *kingpin.Application, i interface{}) string {
	// Prepare
	var depth int
	var fullCommand string
	var help string
	var flags []*kingpin.FlagModel
	var args []*kingpin.ArgModel
	var cmds []*kingpin.CmdModel

	// ApplicationModel vs CmdModel
	switch i.(type) {
	case kingpin.CmdModel:
		cmd := i.(kingpin.CmdModel)
		depth = cmd.Depth
		fullCommand = appCmd.Name + " " + cmd.FullCommand
		// generic
		help = cmd.Help
		if len(cmd.Flags) != 0 {
			flags = cmd.Flags
			fullCommand += " " + cmd.FlagSummary()
		}
		if len(cmd.Args) != 0 {
			args = cmd.Args
			fullCommand += " " + cmd.ArgSummary()
		}
		cmds = cmd.Commands

	case kingpin.ApplicationModel:
		cmd := i.(kingpin.ApplicationModel)
		fullCommand = cmd.Name
		// generic
		help = cmd.Help
		if len(cmd.Flags) != 0 {
			flags = cmd.Flags
			fullCommand += " " + cmd.FlagSummary()
		}
		if len(cmd.Args) != 0 {
			args = cmd.Args
			fullCommand += " " + cmd.ArgSummary()
		}
		cmds = cmd.Commands
	}

	level := depth + 1
	line := fmt.Sprintf("\n<h%d>%s</h%d>", level, html.EscapeString(fullCommand), level)
	line += "\n<pre>" + html.EscapeString(strings.TrimSpace(help)) + "</pre>"
	details := formatFlags(flags) + formatArgs(args)
	if details != "" {
		line += "\n<table>" + details + "\n</table>"
	}
	line += formatCommands(cmds)
	return line
}

func Docs() error {
	m := appCmd.Model()
	result := formatCommand(appCmd, *m)
	fmt.Println(result)
	return nil
}
