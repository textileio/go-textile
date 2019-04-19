package cmd

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
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
Omit the --thread option to use the default thread (if selected).`
}

func (x *chatCmd) Execute(args []string) error {
	setApi(x.Client)

	if x.Thread == "" {
		x.Thread = "default"
	}

	contact, err := getContact()
	if err != nil {
		return err
	}

	rl, err := readline.New(Green(contact.Name + "  "))
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

				btype, err := core.FeedItemType(update)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				if btype != pb.Block_TEXT {
					continue
				}

				payload := new(pb.Text)
				if err := ptypes.UnmarshalAny(update.Payload, payload); err != nil {
					fmt.Println(err.Error())
					continue
				}

				if payload.User.Address != contact.Address {
					if last {
						println()
					}
					println(Cyan(payload.User.Name) + "  " + Grey(payload.Body))
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

func getContact() (*pb.Contact, error) {
	_, c, err := callGetAccountContact()
	if err != nil {
		return nil, err
	}

	if c.Name == "" {
		if len(c.Address) >= 7 {
			c.Name = c.Address[:7]
		} else {
			c.Name = c.Address
		}
	}

	return c, nil
}
