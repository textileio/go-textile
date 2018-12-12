package wallet

import (
	"errors"

	"github.com/textileio/textile-go/keypair"
	"github.com/tyler-smith/go-bip39"
)

var ErrInvalidWordCount = errors.New("invalid word count (must be 12, 15, 18, 21, or 24)")

type WordCount int

const (
	TwelveWords     WordCount = 12
	FifteenWords    WordCount = 15
	EighteenWords   WordCount = 18
	TwentyOneWords  WordCount = 21
	TwentyFourWords WordCount = 24
)

func NewWordCount(cnt int) (*WordCount, error) {
	var wc WordCount
	switch cnt {
	case 12:
		wc = TwelveWords
	case 15:
		wc = FifteenWords
	case 18:
		wc = EighteenWords
	case 21:
		wc = TwentyOneWords
	case 24:
		wc = TwentyFourWords
	default:
		return nil, ErrInvalidWordCount
	}
	return &wc, nil
}

func (w WordCount) EntropySize() int {
	switch w {
	case TwelveWords:
		return 128
	case FifteenWords:
		return 160
	case EighteenWords:
		return 192
	case TwentyOneWords:
		return 224
	case TwentyFourWords:
		return 256
	default:
		return 256
	}
}

// Wallet is a BIP32 Hierarchical Deterministic Wallet based on stellar's
// implementation of https://github.com/satoshilabs/slips/blob/master/slip-0010.md,
// https://github.com/stellar/stellar-protocol/pull/63
type Wallet struct {
	RecoveryPhrase string
}

// NewWallet creates a new wallet with the given entropy size
func NewWallet(entropySize int) (*Wallet, error) {
	entropy, err := bip39.NewEntropy(entropySize)
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
		if err == bip39.ErrInvalidMnemonic {
			return nil, errors.New("invalid mnemonic phrase")
		}
		return nil, err
	}
	masterKey, err := DeriveForPath(TextileAccountPrefix, seed)
	if err != nil {
		return nil, err
	}
	key, err := masterKey.Derive(FirstHardenedIndex + uint32(index))
	if err != nil {
		return nil, err
	}
	return keypair.FromRawSeed(key.RawSeed())
}
