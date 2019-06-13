package cmd

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
)

func Chat(threadID string) error {
	contact, err := getAccountContact()
	if err != nil {
		return err
	}

	rl, err := readline.New(Green(contact.Name + "  "))
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	updates, err := Observe(threadID, []string{"text"})
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

		if err := handleLine(line, threadID); err != nil {
			return err
		}
		last = true
	}
	return nil
}

func handleLine(line string, threadID string) error {
	if strings.TrimSpace(line) != "" {
		if _, err := addMessage(threadID, line); err != nil {
			return err
		}
	}
	return nil
}

func getAccountContact() (*pb.Contact, error) {
	_, c, err := getAccount()
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
