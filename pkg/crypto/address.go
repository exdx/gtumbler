package crypto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"net/http"
)

const (
	sendURL = "http://jobcoin.gemini.com/survey/api/transactions"
	checkURL = "http://jobcoin.gemini.com/survey/api/addresses/"
)

type Address string
type Amount string

// CreateAddress generates an address - the address is a valid ethereum address
// The private key is discarded as it is not required
func CreateAddress() (Address, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}

	// address resembles something like 0x8ee3333cDE801ceE9471ADf23370c48b011f82a6
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()

	return Address(address), nil
}

// Send physically sends coins from "from" to "to" over the protocol
func Send(from Address, to Address, size Amount) error {
	request := &SendCoinRequest{
		From: from,
		To: to,
		Amount: size,
	}

	req, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(sendURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 422 {
		return fmt.Errorf("insufficient funds in address %s", from)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unknown error when sending to deposit address")
	}

	return nil
}

// CheckAddress checks the balance of an address and returns the number of coins held, which is at-least 0
// Sent over the protocol
func CheckAddress(address Address) (Amount, error) {
	target := fmt.Sprint(checkURL, address)
	result := &CheckAddressResponse{}

	resp, err := http.Get(target)
	if err != nil {
		return Amount(0), err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Amount(0), err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return Amount(0), err
	}

	return result.Balance, nil
}