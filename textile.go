package main

import (
	"fmt"
	"os"

	"github.com/textileio/go-textile/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
