package cmd

import (
	"strings"

	"github.com/chzyer/readline"
)

func init() {
	register(&chatCmd{})
}

type chatCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
}

func (x *chatCmd) Name() string {
	return "chat"
}

func (x *chatCmd) Short() string {
	return "Start a thread chat"
}

func (x *chatCmd) Long() string {
	return `
Starts an interactive chat session in a thread.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *chatCmd) Execute(args []string) error {
	setApi(x.Client)

	if x.Thread == "" {
		x.Thread = "default"
	}

	pid, username, err := getUsername()
	if err != nil {
		return err
	}

	rl, err := readline.New(Green(username + "  "))
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	updates, err := callSub(x.Thread, []string{"message"})
	if err != nil {
		return err
	}

	last := true
	go func() {
		for {
			select {
			case update, ok := <-updates:
				if !ok {
					return
				}
				if update.Block.AuthorId != pid {
					if last {
						println()
					}
					println(Cyan(update.Block.Username) + "  " + Grey(update.Block.Body))
					last = false
				}
			}
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		if err := handleLine(line, x.Thread); err != nil {
			return err
		}
		last = true
	}
	return nil
}

func handleLine(line string, threadId string) error {
	if strings.TrimSpace(line) != "" {
		if _, err := callAddMessages(threadId, line); err != nil {
			return err
		}
	}
	return nil
}

func getUsername() (string, string, error) {
	_, prof, err := callGetProfile()
	if err != nil {
		return "", "", err
	}
	username := prof.Username

	pid, err := callId()
	if err != nil {
		return "", "", err
	}
	if username == "" {
		if len(pid) >= 7 {
			username = pid[len(pid)-7:]
		} else {
			username = pid
		}
	}

	return pid, username, nil
}
