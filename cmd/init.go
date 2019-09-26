package cmd

import (
	"fmt"

	"github.com/textileio/go-textile/core"
)

func InitCommand(config core.InitConfig) error {
	if err := core.InitRepo(config); err != nil {
		return fmt.Errorf("initialize failed: %s", err)
	}
	fmt.Printf("Initialized account with address %s\n", config.Account.Address())
	return nil
}
