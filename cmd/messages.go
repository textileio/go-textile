package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

var errMissingMessageBody = errors.New("missing message body")
var errMissingMessageId = errors.New("missing message block ID")

func init() {
	register(&messagesCmd{})
}

type messagesCmd struct {
	Add    addMessagesCmd `command:"add" description:"Add a thread message"`
	List   lsMessagesCmd  `command:"ls" description:"List thread messages"`
	Get    getMessagesCmd `command:"get" description:"Get a thread message"`
	Ignore rmMessagesCmd  `command:"ignore" description:"Ignore a thread message"`
}

func (x *messagesCmd) Name() string {
	return "messages"
}

func (x *messagesCmd) Short() string {
	return "Manage thread messages"
}

func (x *messagesCmd) Long() string {
	return `
Messages are added as blocks in a thread.
Use this command to add, list, get, and ignore messages.
`
}

type addMessagesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
}

func (x *addMessagesCmd) Usage() string {
	return `

Adds a message to a thread.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *addMessagesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingMessageBody
	}

	if x.Thread == "" {
		x.Thread = "default"
	}

	res, err := callAddMessages(x.Thread, args[0])
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func callAddMessages(threadId string, body string) (string, error) {
	res, err := executeJsonCmd(POST, "threads/"+threadId+"/messages", params{
		args: []string{body},
	}, nil)
	if err != nil {
		return "", err
	}
	return res, nil
}

type lsMessagesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  int           `short:"l" long:"limit" description:"List page size." default:"10"`
}

func (x *lsMessagesCmd) Usage() string {
	return `

Paginates thread messages.
Omit the --thread option to paginate all messages.
Specify "default" to use the default thread (if selected).
`
}

func (x *lsMessagesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"offset": x.Offset,
		"limit":  strconv.Itoa(x.Limit),
	}
	return callLsMessages(opts)
}

func callLsMessages(opts map[string]string) error {
	var list pb.TextList
	res, err := executeJsonPbCmd(GET, "messages", params{opts: opts}, &list)
	if err != nil {
		return err
	}
	if len(list.Items) > 0 {
		output(res)
	}

	limit, err := strconv.Atoi(opts["limit"])
	if err != nil {
		return err
	}
	if len(list.Items) < limit {
		return nil
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return callLsMessages(map[string]string{
		"thread": opts["thread"],
		"offset": list.Items[len(list.Items)-1].Block,
		"limit":  opts["limit"],
	})
}

type getMessagesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getMessagesCmd) Usage() string {
	return `

Gets a thread message by block ID.`
}

func (x *getMessagesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingMessageId
	}

	res, err := executeJsonCmd(GET, "messages/"+util.TrimQuotes(args[0]), params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type rmMessagesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmMessagesCmd) Usage() string {
	return `

Ignores a thread message by its block ID.
This adds an "ignore" thread block targeted at the message.
Ignored blocks are by default not returned when listing. 
`
}

func (x *rmMessagesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmBlocks(args)
}
