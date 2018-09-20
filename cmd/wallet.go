package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/util"
	"github.com/tyler-smith/go-bip39"
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

	var wcount util.WordCount
	count := c.MultiChoice([]string{
		"12",
		"15",
		"18",
		"21",
		"24",
	}, "How many words?")

	switch count {
	case 0:
		wcount = util.TwelveWords
	case 1:
		wcount = util.FifteenWords
	case 2:
		wcount = util.EighteenWords
	case 3:
		wcount = util.TwentyOneWords
	case 4:
		wcount = util.TwentyFourWords
	default:
		c.Err(errors.New("invalid word count"))
	}

	mnemonic, err := util.CreateMnemonic(wcount)
	if err != nil {
		c.Err(err)
		return
	}

	c.Println(cyan(strings.Repeat("-", len(mnemonic)+4)))
	c.Println(cyan("| " + mnemonic + " |"))
	c.Println(cyan(strings.Repeat("-", len(mnemonic)+4)))
	c.Println(grey("WARNING! Store these words above in a safe place!"))
	c.Println(grey("WARNING! If you lose your words, you will lose access to data in all derived accounts!"))
	c.Println(grey("WARNING! Anyone who has access to these words can access your wallet accounts!"))
	c.Println("")
	c.Println(grey("Use: `wallet accounts` command to see all generated accounts."))
	c.Println("")

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		c.Err(err)
		return
	}

	masterKey, err := crypto.DeriveForPath(crypto.TextileAccountPrefix, seed)
	if err != nil {
		c.Err(err)
		return
	}

	key, err := masterKey.Derive(crypto.FirstHardenedIndex)
	if err != nil {
		c.Err(err)
		return
	}
	kp, err := keypair.FromRawSeed(key.RawSeed())
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

	mnemonic := strings.Join(words, " ")
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		c.Err(errors.New("invalid words or checksum"))
		return
	}

	masterKey, err := crypto.DeriveForPath(crypto.TextileAccountPrefix, seed)
	if err != nil {
		c.Err(err)
		return
	}

	i := 0
	more := true
	for more {
		key, err := masterKey.Derive(crypto.FirstHardenedIndex + uint32(i))
		if err != nil {
			c.Err(err)
			return
		}
		kp, err := keypair.FromRawSeed(key.RawSeed())
		if err != nil {
			c.Err(err)
			return
		}
		c.Println(grey(fmt.Sprintf("--- ACCOUNT %d ---", i)))
		c.Println(green(fmt.Sprintf("PUBLIC KEY: %s", kp.Address())))
		c.Println(red(fmt.Sprintf("SECRET KEY: %s", kp.Seed())))

		c.Print("See next account (Y/n)? ")
		ans := c.ReadLine()
		if ans == "n" {
			more = false
		}
		i++
	}
}
