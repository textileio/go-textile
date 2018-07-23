package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func CafeReferral(c *ishell.Context) {
	c.Print("key: ")
	password := c.ReadPassword()

	if err := core.Node.Wallet.GetReferral(creds); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("welcome back, %s!", username)))
}

func CafeRegister(c *ishell.Context) {
	c.Print("email address: ")
	email := c.ReadLine()
	c.Print("username: ")
	username := c.ReadLine()
	c.Print("referral code: ")
	code := c.ReadLine()
	c.Print("password: ")
	password := c.ReadPassword()

	reg := &models.Registration{
		Username: username,
		Password: password,
		Identity: &models.Identity{
			Type:  models.EmailAddress,
			Value: email,
		},
		Referral: code,
	}
	if err := core.Node.Wallet.SignUp(reg); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("welcome aboard, %s!", username)))
}

func CafeLogin(c *ishell.Context) {
	c.Print("username: ")
	username := c.ReadLine()
	c.Print("password: ")
	password := c.ReadPassword()

	creds := &models.Credentials{
		Username: username,
		Password: password,
	}
	if err := core.Node.Wallet.SignIn(creds); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("welcome back, %s!", username)))
}

func CafeLogout(c *ishell.Context) {
	c.Print("logout? Y/n")
	confirm := c.ReadLine()

	if confirm != "" && confirm != "Y" {
		return
	}
	if err := core.Node.Wallet.SignOut(); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("see ya!"))
}

func CafeStatus(c *ishell.Context) {
	signedIn, err := core.Node.Wallet.IsSignedIn()
	if err != nil {
		c.Err(err)
		return
	}
	if signedIn {
		c.Println(color.New(color.FgHiGreen).SprintFunc()("logged in"))
	} else {
		c.Println(color.New(color.FgHiRed).SprintFunc()("not logged in"))
	}
}
