package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/wallet"
	"gopkg.in/abiosoft/ishell.v2"
	"regexp"
	"strings"
)

var wordsRegexp = regexp.MustCompile(`^[a-z]+$`)

func CreateWallet(c *ishell.Context) {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	grey := color.New(color.FgHiBlack).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()

	var wcount wallet.WordCount
	count := c.MultiChoice([]string{
		"12",
		"15",
		"18",
		"21",
		"24",
	}, "How many words?")

	switch count {
	case 0:
		wcount = wallet.TwelveWords
	case 1:
		wcount = wallet.FifteenWords
	case 2:
		wcount = wallet.EighteenWords
	case 3:
		wcount = wallet.TwentyOneWords
	case 4:
		wcount = wallet.TwentyFourWords
	default:
		c.Err(errors.New("invalid word count"))
	}

	w, err := wallet.NewWallet(wcount)
	if err != nil {
		c.Err(err)
		return
	}

	c.Println(cyan(strings.Repeat("-", len(w.RecoveryPhrase)+4)))
	c.Println(cyan("| " + w.RecoveryPhrase + " |"))
	c.Println(cyan(strings.Repeat("-", len(w.RecoveryPhrase)+4)))
	c.Println(grey("WARNING! Store these words above in a safe place!"))
	c.Println(grey("WARNING! If you lose your words, you will lose access to data in all derived accounts!"))
	c.Println(grey("WARNING! Anyone who has access to these words can access your wallet accounts!"))
	c.Println("")
	c.Println(grey("Use: `wallet accounts` command to see all generated accounts."))
	c.Println("")

	c.Print("Enter password (leave empty if none): ")
	password := c.ReadLine()

	kp, err := w.AccountAt(0, password)
	if err != nil {
		c.Err(err)
		return
	}

	c.Println(grey("--- ACCOUNT 0 ---"))
	c.Println(green(fmt.Sprintf("PUBLIC KEY: %s", kp.Address())))
	c.Println(red(fmt.Sprintf("SECRET KEY: %s", kp.Seed())))
}

func WalletAccounts(c *ishell.Context) {
	grey := color.New(color.FgHiBlack).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()

	var wcount int
	count := c.MultiChoice([]string{
		"12",
		"15",
		"18",
		"21",
		"24",
	}, "How many words are in your mnemonic phrase?")

	switch count {
	case 0:
		wcount = 12
	case 1:
		wcount = 15
	case 2:
		wcount = 18
	case 3:
		wcount = 21
	case 4:
		wcount = 24
	default:
		c.Err(errors.New("invalid word count"))
	}

	words := make([]string, int(wcount))
	for i := 0; i < int(wcount); i++ {
		c.Print(fmt.Sprintf("Enter word #%d: ", i+1))
		words[i] = c.ReadLine()
		if !wordsRegexp.MatchString(words[i]) {
			c.Println("Invalid word, try again.")
			i--
		}
	}
	wall := wallet.NewWalletFromRecoveryPhrase(strings.Join(words, " "))

	c.Print("Enter password (leave empty if none): ")
	password := c.ReadLine()

	i := 0
	more := true
	for more {
		kp, err := wall.AccountAt(i, password)
		if err != nil {
			c.Err(err)
			return
		}
		c.Println(grey(fmt.Sprintf("--- ACCOUNT %d ---", i)))
		c.Println(green(fmt.Sprintf("PUBLIC KEY: %s", kp.Address())))
		c.Println(red(fmt.Sprintf("SECRET KEY: %s", kp.Seed())))

		c.Print("Show next account (Y/n)? ")
		ans := c.ReadLine()
		if ans == "n" {
			more = false
		}
		i++
	}
}
