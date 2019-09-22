package client

import "github.com/Denton24646/gtumbler/pkg/crypto"

type Config struct {
	MixerURL        string         `cfgDefault:"http://localhost:8989/create"`
	NumberAddresses int            `cfgDefault:"3"`
	SendAddress     crypto.Address `cfgDefault:"Genesis"`
	Size            crypto.Amount  `cfgDefault:"4"`
}
