package cmd

import (
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func cafeRegister(c *ishell.Context) {
	c.Print("cafe peer id: ")
	peerId := c.ReadLine()

	if err := core.Node.CafeRegister(peerId); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("welcome!"))
}
