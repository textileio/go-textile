package cmd

import (
	"fmt"
	"strings"

	"github.com/textileio/go-textile/wallet"
)

func WalletInit(words int, passphrase string) error {
	wordCount, err := wallet.NewWordCount(words)
	if err != nil {
		return err
	}

	w, err := wallet.WalletFromEntropy(wordCount.EntropySize())
	if err != nil {
		return err
	}

	// Print the recovery phrase surrounded by a box of dashes
	fmt.Println(strings.Repeat("-", len(w.RecoveryPhrase)+4))
	fmt.Println("| " + w.RecoveryPhrase + " |")
	fmt.Println(strings.Repeat("-", len(w.RecoveryPhrase)+4))
	fmt.Println("WARNING! Store these words above in a safe place!")
	fmt.Println("WARNING! If you lose your words, you will lose access to data in all derived accounts!")
	fmt.Println("WARNING! Anyone who has access to these words can access your wallet accounts!")
	fmt.Println("")
	fmt.Println("Use: `wallet accounts` command to inspect more accounts.")
	fmt.Println("")

	// show first account
	kp, err := w.AccountAt(0, passphrase)
	if err != nil {
		return err
	}
	fmt.Println("--- ACCOUNT 0 ---")
	fmt.Println(kp.Address())
	fmt.Println(kp.Seed())

	return nil
}

func WalletAccounts(mnemonic string, passphrase string, depth int, offset int) error {
	if depth < 1 || depth > 100 {
		return fmt.Errorf("depth must be greater than 0 and less than 100")
	}
	if offset < 0 || offset > depth {
		return fmt.Errorf("offset must be greater than 0 and less than depth")
	}

	wall := wallet.WalletFromMnemonic(mnemonic)

	for i := offset; i < offset+depth; i++ {
		kp, err := wall.AccountAt(i, passphrase)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("--- ACCOUNT %d ---", i))
		fmt.Println(kp.Address())
		fmt.Println(kp.Seed())
	}
	return nil
}
