package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func cafeRegister(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer id"))
		return
	}
	peerId := c.Args[0]

	if err := core.Node.RegisterCafe(peerId); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("welcome!"))
}

func cafeList(c *ishell.Context) {
	cafes, err := core.Node.ListCafeSessions()
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	for _, cafe := range cafes {
		c.Println(green(fmt.Sprintf("peer id: %s, expires: %s", cafe.CafeId, cafe.Expiry.String())))
	}
}

func cafeDeregister(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer id"))
		return
	}
	peerId := c.Args[0]

	if err := core.Node.DeregisterCafe(peerId); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("see ya!"))
}

func cafeCheckMessages(c *ishell.Context) {
	if err := core.Node.CheckCafeMessages(); err != nil {
		c.Println(fmt.Errorf("check messages failed: %s", err))
		return
	}
	c.Println("ok, checking")
}
