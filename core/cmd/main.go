package cmd

type TextResponse struct {
	AccountAddress string `json:"account_address"`
	AccountId      string `json:"account_id"`
	PeerId         string `json:"peer_id"`
}
