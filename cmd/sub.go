package cmd

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/textileio/textile-go/core"
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

-  JOIN
-  ANNOUNCE
-  LEAVE
-  MESSAGE
-  FILES
-  COMMENT
-  LIKE
-  MERGE
-  IGNORE
-  FLAG

Use the --thread option to subscribe to events emmited from a specific thread.

Use the --type option to limit the output to specific update type(s).
This option can be used multiple times, e.g., --type files --type comment.
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

			data, err := json.MarshalIndent(update, "", "    ")
			if err == io.EOF {
				break
			} else if err != nil {
				return nil
			}
			output(string(data))
		}
	}
}

func callSub(threadId string, types []string) (<-chan core.ThreadUpdate, error) {
	if threadId != "" {
		threadId = "/" + threadId
	}

	updates := make(chan core.ThreadUpdate, 10)
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
			var info core.ThreadUpdate
			if err := decoder.Decode(&info); err == io.EOF {
				return
			} else if err != nil {
				output(err.Error())
				return
			}
			updates <- info
		}
	}()

	return updates, nil
}
