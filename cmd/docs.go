

package cmd

import (
	"fmt"
	"strings"
	"html"
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
		line += "\n<table>\n\t<tr><th>Flag</th><th>Description</th></tr>"
		for _, flag := range flags {
			name := formatFlag(flag.Short != 0, flag)
			desc := strings.TrimSpace(flag.Help)
			line += "\n\t<tr><td>" + html.EscapeString(name) + "</td><td><pre>" + html.EscapeString(desc) + "</pre></td></tr>"
		}
		line += "\n</table>"
	}
	return line
}

func formatArgs(args []*kingpin.ArgModel) string {
	var line = ""
	if len(args) != 0 {
		line += "\n<table>\n\t<tr><th>Argument</th><th>Description</th></tr>"
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
			line += "\n\t<tr><td>" + html.EscapeString(name) + "</td><td><pre>" + html.EscapeString(desc) + "</pre></td></tr>"
		}
		line += "\n</table>"
	}
	return line
}

func formatCommands(cmds []*kingpin.CmdModel) string {
	var line = ""
	if len(cmds) != 0 {
		for _, subCmd := range cmds {
			line += "\n" + formatCommand(appCmd, subCmd)
		}
	}
	return line
}

func formatAppCommand (cmd *kingpin.ApplicationModel) string {
	line := appCmd.Name
	if len(cmd.Flags) != 0 {
		line += " " + cmd.FlagSummary()
	}
	if len(cmd.Args) != 0 {
		line += " " + cmd.ArgSummary()
	}
	line = fmt.Sprintf("\n<h%d>%s</h%d>", 1, html.EscapeString(line), 1)
	desc := strings.TrimSpace(cmd.Help)
	line += "\n<p><pre>" + html.EscapeString(desc) + "</pre></p>"
	line += formatFlags(cmd.Flags)
	line += formatArgs(cmd.Args)
	line += formatCommands(cmd.Commands)
	return line
}

func formatCommand(appCmd *kingpin.Application, cmd *kingpin.CmdModel) string {
	level := cmd.Depth + 1
	line := appCmd.Name + " " + cmd.FullCommand
	if len(cmd.Flags) != 0 {
		line += " " + cmd.FlagSummary()
	}
	if len(cmd.Args) != 0 {
		line += " " + cmd.ArgSummary()
	}
	line = fmt.Sprintf("\n<h%d>%s</h%d>", level, html.EscapeString(line), level)
	desc := strings.TrimSpace(cmd.Help)
	line += "\n<p><pre>" + html.EscapeString(desc) + "</pre></p>"
	line += formatFlags(cmd.Flags)
	line += formatArgs(cmd.Args)
	line += formatCommands(cmd.Commands)
	return line
}

func Docs() error {
	m := appCmd.Model()
	result := formatAppCommand(m)
	fmt.Println(result)
	return nil
}