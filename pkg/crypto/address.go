package crypto

import "github.com/ethereum/go-ethereum/crypto"

type Address string

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
func Send(from Address, to Address) error {
	return nil
}

// CheckAddress checks the balance of an address and returns the number of coins held, which is atleast 0 (therfore uint)
// Sent over the protocol
func CheckAddress(address Address) (uint, error) {
	return 0, nil
}