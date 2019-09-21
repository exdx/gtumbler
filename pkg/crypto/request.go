package crypto

type SendCoinRequest struct {
	From   Address `json:"fromAddress"`
	To     Address `json:"toAddress"`
	Amount Amount  `json:"amount"`
}

// TODO parse additional address information: transaction array
// TODO for now we are only interested in the total balance of the address being greater than zero, assuming its new
type CheckAddressResponse struct {
	Balance Amount `json:"balance"`
}