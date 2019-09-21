package tumbler

import "github.com/Denton24646/gtumbler/pkg/crypto"

// tumbler is responsible for the actual mixing process
// given an array of addresses and a deposit address it must do the following
// 1. split the amount in the deposit address into random sizes
// 2. allocate the random sizes into house array of addresses at random times
// 3. check if others are also mixing and mix their coins in as well
// 3. report the initial tumbling step as complete

// for sending the coins back out to the customer the tumbler must do the opposite
// 1. given an array of house addresses it must send the appropriate amount of coins to the array of customer accounts
// 2. it must send them in random sizes over random periods of time
// 3. report the final tumbling process as complete

type Tumble interface {
	// Mix mixes the client coins from the deposit address back to various house addresses
	Mix(depositAddress crypto.Address, houseAddresses []crypto.Address) error
	//SendMixedFunds sends the funds back to the customer specified accounts from house addresses
	SendMixedFunds(customerAddresses []crypto.Address, houseAddresses []crypto.Address) error
}

type Tumbler struct {}

func New() *Tumbler {
	return &Tumbler{}
}

func (t *Tumbler) Mix(depositAddress crypto.Address, houseAddresses []crypto.Address) error {
	return nil
}

func (t *Tumbler) SendMixedFunds(customerAddresses []crypto.Address, houseAddresses []crypto.Address) error {
	return nil
}