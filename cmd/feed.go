package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

func Feed(threadID string, offset string, limit int, mode string) error {
	var list pb.FeedItemList
	opts := map[string]string{
		"thread": threadID,
		"offset": offset,
		"limit":  strconv.Itoa(limit),
		"mode":   mode,
	}
	res, err := executeJsonPbCmd(http.MethodGet, "feed", params{opts: opts}, &list)
	if err != nil {
		return err
	}
	if list.Count > 0 {
		output(res)
	}

	if list.Next == "" {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return Feed(threadID, list.Next, limit, mode)
}
