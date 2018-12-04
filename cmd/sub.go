package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/textileio/textile-go/core"
)

func init() {
	register(&subCmd{})
}

type subCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Type   string        `short:"k" long:"type" description:"Comma-separated list of event types to include. Omit for all."`
}

func (x *subCmd) Name() string {
	return "sub"
}

func (x *subCmd) Short() string {
	return "Subscribe to thread events/updates"
}

func (x *subCmd) Long() string {
	return `
Subscribe to thread events/updates.
Use the --thread option to subscribe to events emmited from a specific thread.  
Use the --type option to filter to specific event type(s),
e.g., files,comments,likes. Omit for all.
`
}

func (x *subCmd) Execute(args []string) error {
	setApi(x.Client)

	if x.Thread != "" {
		x.Thread = "/" + x.Thread
	}

	// '|' doesn't work on cmdline, so use commas (',') and swap out for '|'
	x.Type = strings.Join(strings.Split(x.Type, ","), "|")

	req, err := request(GET, "sub"+x.Thread, params{
		opts: map[string]string{
			"thread": x.Thread,
			"type":   x.Type,
		},
	})
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode >= 400 {
		res, err := unmarshalString(req.Body)
		if err != nil {
			return err
		}
		return errors.New(res)
	}
	decoder := json.NewDecoder(req.Body)
	for {
		var info core.ThreadUpdate
		if err := decoder.Decode(&info); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		data, err := json.MarshalIndent(info, "", "    ")
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		output(string(data), nil)
	}
	return nil
}
