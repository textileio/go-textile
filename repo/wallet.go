package repo

import "time"

type Photo map[string]string

type WalletData struct {
	Photos []Photo `json:"photos"`
}

type Wallet struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Data WalletData `json:"data"`
}

//func (w *Wallet) PinPhoto() (string, error) {
//
//}
