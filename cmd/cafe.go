package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"strconv"
)

func CafeReferral(c *ishell.Context) {
	c.Print("key: ")
	key := c.ReadPassword()
	c.Print("count (1): ")
	counts := c.ReadLine()
	c.Print("limit (1): ")
	limits := c.ReadLine()

	count, err := strconv.Atoi(counts)
	if err != nil {
		count = 1
	}
	limit, err := strconv.Atoi(limits)
	if err != nil {
		limit = 1
	}
	username, err := core.Node.Wallet.GetUsername()
	if err != nil {
		c.Err(err)
		return
	}
	if username == nil {
		tmp := "anonymous"
		username = &tmp
	}
	req := &models.ReferralRequest{
		Key:         key,
		Count:       count,
		Limit:       limit,
		RequestedBy: *username,
	}
	res, err := core.Node.Wallet.CreateCafeReferral(req)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	for _, ref := range res.RefCodes {
		c.Println(green(ref))
	}
}

func ListCafeReferrals(c *ishell.Context) {
	c.Print("key: ")
	key := c.ReadPassword()

	res, err := core.Node.Wallet.ListCafeReferrals(key)
	if err != nil {
		c.Err(err)
		return
	}

	yellow := color.New(color.FgHiYellow).SprintFunc()
	for _, ref := range res.RefCodes {
		c.Println(yellow(ref))
	}
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

	reg := &models.UserRegistration{
		Username: username,
		Password: password,
		Identity: &models.UserIdentity{
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

	creds := &models.UserCredentials{
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

func CafeTokens(c *ishell.Context) {
	tokens, err := core.Node.Wallet.GetCafeTokens(false)
	if err != nil {
		c.Err(err)
		return
	}
	if tokens == nil {
		c.Println(color.New(color.FgHiRed).SprintFunc()("no tokens found"))
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("expiry: %s", tokens.Expiry.String())))
}
