package tumbler

import (
	"fmt"
	"github.com/Denton24646/gtumbler/pkg/crypto"
	"strconv"
)

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

const minDeposit = 0.1
const maxDeposit = 100

type Tumble interface {
	// Mix mixes the client coins from the deposit address back to various house addresses
	Mix(depositAddress crypto.Address, houseAddresses []crypto.Address) error
	//SendMixedFunds sends the funds back to the customer specified accounts from house addresses
	SendMixedFunds(customerAddresses []crypto.Address, houseAddresses []crypto.Address) error
}

type Tumbler struct {
	Size crypto.Amount
	Strategies strategies
}

func New(amount crypto.Amount) *Tumbler {
	return &Tumbler{
		Size: amount,
		Strategies: getStrategies(),
	}
}

// Mix mixes coins on the front-end of the transaction, from the deposit address to the house
// It has information from the mixer about how many coins there are deposited
// Then it uses some randomness to send those funds along to random houseAddresses
func (t *Tumbler) Mix(depositAddress crypto.Address, houseAddresses []crypto.Address) error {
	//parse amount, which is provided as a string, into a float
	amount, err := strconv.ParseFloat(string(t.Size), 64)
	if err != nil {
		return err
	}
	// validate amount deposited is valid
	if !valid(amount) {
		return fmt.Errorf("funds are not within the specified guidelines for gtumbler")
	}

	// pick random strategy from map
	strategyKey := pickRandom(len(t.Strategies))
	strategy := t.Strategies[strategyKey]

	// send amount in strategy to random house address
	// TODO use some time variability to add additional randomness
	var sendAmount float64
	for _, chunk := range strategy {
		amount * chunk = sendAmount
		houseKey := pickRandom(len(houseAddresses))
		err := crypto.Send(depositAddress, houseAddresses[houseKey], crypto.Amount(fmt.Sprintf("%f", sendAmount)))
		if err != nil {
			return err
		}
	}

	return nil
}

// Deposits need to be validated: they have a certain minimum and maximum size
// This is to ensure the mixer has enough liquidity to mix all customer deposits
func valid(size float64) bool {
	if size > minDeposit && size < maxDeposit {
		return true
	}
	return false
}

// SendMixedFunds sends funds from random house addresses to the customer deposit addresses
func (t *Tumbler) SendMixedFunds(customerAddresses []crypto.Address, houseAddresses []crypto.Address) error {
	//parse amount, which is provided as a string, into a float
	amount, err := strconv.ParseFloat(string(t.Size), 64)
	if err != nil {
		return err
	}

	// pick random strategy from map
	strategyKey := pickRandom(len(t.Strategies))
	strategy := t.Strategies[strategyKey]

	// send funds from a random house address to a random customer address
	// note: this does not ensure each address the customer specified will receive funds, for example one may receive all funds
	var sendAmount float64
	for _, chunk := range strategy {
		amount * chunk = sendAmount
		houseKey := pickRandom(len(houseAddresses))
		customerKey := pickRandom(len(customerAddresses))
		err := crypto.Send(houseAddresses[houseKey], customerAddresses[customerKey],  crypto.Amount(fmt.Sprintf("%f", sendAmount)))
		if err != nil {
			return err
		}
	}

	return nil
}