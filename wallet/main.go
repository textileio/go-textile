package wallet

import (
	"github.com/textileio/textile-go/keypair"
	"github.com/tyler-smith/go-bip39"
)

type WordCount int

const (
	TwelveWords     WordCount = 128
	FifteenWords              = 160
	EighteenWords             = 192
	TwentyOneWords            = 224
	TwentyFourWords           = 256
)

// Wallet is a BIP32 Hierarchical Deterministic Wallet based on stellar's
// implementation of https://github.com/satoshilabs/slips/blob/master/slip-0010.md,
// https://github.com/stellar/stellar-protocol/pull/63
type Wallet struct {
	RecoveryPhrase string
}

// NewWallet creates a new wallet with a mnemonic recovery phrase and
func NewWallet(wc WordCount) (*Wallet, error) {
	entropy, err := bip39.NewEntropy(int(wc))
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return &Wallet{RecoveryPhrase: mnemonic}, nil
}

func NewWalletFromRecoveryPhrase(mnemonic string) *Wallet {
	return &Wallet{RecoveryPhrase: mnemonic}
}

func (w *Wallet) AccountAt(index int, password string) (*keypair.Full, error) {
	seed, err := bip39.NewSeedWithErrorChecking(w.RecoveryPhrase, password)
	if err != nil {
		return nil, err
	}
	masterKey, err := DeriveForPath(TextileAccountPrefix, seed)
	if err != nil {
		return nil, err
	}
	key, err := masterKey.Derive(FirstHardenedIndex)
	if err != nil {
		return nil, err
	}
	return keypair.FromRawSeed(key.RawSeed())
}
