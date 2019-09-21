package client

import "github.com/Denton24646/gtumbler/pkg/crypto"

type Config struct {
	MixerURL        string //default "localhost:8989/create"
	NumberAddresses int    // number of addresses, default 5
	SendAddress     crypto.Address
	Size            crypto.Amount
}
