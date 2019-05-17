package main

import (
	"fmt"
	"github.com/textileio/go-textile/cmd"
	"os"
)

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
