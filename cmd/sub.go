package cmd

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/util"
)

func init() {
	register(&subCmd{})
}

type subCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Type   []string      `short:"k" long:"type" description:"An update type to filter for. Omit for all."`
}

func (x *subCmd) Name() string {
	return "sub"
}

func (x *subCmd) Short() string {
	return "Subscribe to thread updates"
}

func (x *subCmd) Long() string {
	return `
Subscribes to updates in a thread or all threads. An update is generated
when a new block is added to a thread.

There are several update types:

-  MERGE
-  IGNORE
-  FLAG
-  JOIN
-  ANNOUNCE
-  LEAVE
-  TEXT
-  FILES
-  COMMENT
-  LIKE

Use the --thread option to subscribe to events emmited from a specific thread.
The --type option can be used multiple times, e.g., --type files --type comment.
`
}

func (x *subCmd) Execute(args []string) error {
	setApi(x.Client)

	updates, err := callSub(x.Thread, x.Type)
	if err != nil {
		return err
	}

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return nil
			}

			out, err := pbMarshaler.MarshalToString(update)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			output(out)
		}
	}
}

func callSub(threadId string, types []string) (<-chan *pb.FeedItem, error) {
	if threadId != "" {
		threadId = "/" + threadId
	}

	updates := make(chan *pb.FeedItem, 10)
	go func() {
		defer close(updates)

		res, cancel, err := request(GET, "sub"+threadId, params{
			opts: map[string]string{"type": strings.Join(types, "|")},
		})
		if err != nil {
			output(err.Error())
			return
		}
		defer res.Body.Close()
		defer cancel()

		if res.StatusCode >= 400 {
			body, err := util.UnmarshalString(res.Body)
			if err != nil {
				output(err.Error())
			} else {
				output(body)
			}
			return
		}

		decoder := json.NewDecoder(res.Body)
		for decoder.More() {
			var update pb.FeedItem
			if err := pbUnmarshaler.UnmarshalNext(decoder, &update); err == io.EOF {
				return
			} else if err != nil {
				output(err.Error())
				return
			}
			updates <- &update
		}
	}()

	return updates, nil
}
