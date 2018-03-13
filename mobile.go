package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/textileio/textile-go/mobile"
	"os"
)

var opts struct {
	RepoDir string `short:"d" long:"dir" description:"Rep directory" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	textile := mobile.NewTextile(opts.RepoDir)
	textile.Start()
}
