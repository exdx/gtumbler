package models

import "github.com/Denton24646/gtumbler/pkg/crypto"

type CleanAddressRequest struct {
	Id        int              `json:"id"`
	Addresses []crypto.Address `json:"addresses"`
}

type CleanAddressResponse struct {
	DepositAddress crypto.Address `json:"address"`
}


