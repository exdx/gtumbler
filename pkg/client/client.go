package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Denton24646/gtumbler/pkg/crypto"
	"github.com/Denton24646/gtumbler/pkg/models"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// The client is an http client that sends requests to the mixer consisting of the following
// 1. The initial request sends a list of new addresses that the mixed coins will eventually be sent back to
// 2. The mixer responds with a deposit address
// 3. The client sends the full deposit amount to the deposit address
// 4. From that point on the client checks the list of addresses sent in (1) to be notified when their mixing coins are available

type Client interface {
	CreateCleanAddresses(number int) ([]crypto.Address, error)
	SendCleanAddresses() error
	SendDeposit(address crypto.Address) error
	CheckCleanAddresses() (bool, error)
}

type UserClient struct {
	// Id is pseudo-random id shared between client and server
	Id int
	// mixerURL is the location of the mixer server (localhost:8989 when running locally)
	mixerURL string
	// List of clean addresses the client wants the coins to end up in: these can be generated or provided at runtime
	CleanAddresses []crypto.Address
	// Deposit address that the user client receives from the server
	DepositAddress crypto.Address
	// Timestamp of when the client deposit was sent
	SentTimestamp time.Time
	// Timestamp of when the funds were deposited across all customer addresses (the mixing is complete)
	ReceivedTimestamp time.Time
}

func New(config Config) *UserClient {
	return &UserClient{
		Id:       rand.Int(),
		mixerURL: config.MixerURL,
	}
}

// CreateCleanAddresses generates the number of addresses specified by the caller
// If there is an error, terminate and return an empty list - address creation is atomic: either they are all created or it fails
func (u *UserClient) CreateCleanAddresses(number int) ([]crypto.Address, error) {
	var addresses []crypto.Address
	for i := 0; i < number; i++ {
		a, err := crypto.CreateAddress()
		if err != nil {
			return nil, err
		}
		fmt.Printf("  Address %d: %s\n", i, a)
		addresses = append(addresses, a)
	}
	u.CleanAddresses = addresses
	return addresses, nil
}

// SendCleanAddresses sends the clean addresses to the mixer in an http POST request to the specified endpoint
// The mixer sends the deposit address in the response to the request
func (u *UserClient) SendCleanAddresses() error {
	request := models.CleanAddressRequest{
		Id:        u.Id,
		Addresses: u.CleanAddresses,
	}

	req, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(u.mixerURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := &models.CleanAddressResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return err
	}

	u.DepositAddress = response.DepositAddress
	return nil
}

// SendDeposit sends coins to the deposit address specified by the mixer from an arbitrary address
func (u *UserClient) SendDeposit(address crypto.Address, size crypto.Amount) error {
	err := crypto.Send(address, u.DepositAddress, size)
	if err != nil {
		return err
	}
	return nil
}

// CheckCleanAddresses checks to see if all the provided addresses received the end deposits from the mixer
// It returns true in the case where at least one provided address has funds, false otherwise
// Assumes provided addresses have a zero balance initially
func (u *UserClient) CheckCleanAddresses() (bool, error) {
	var found bool
	for _, address := range u.CleanAddresses {
		amount, err := crypto.CheckAddress(address)
		if err != nil {
			return false, err
		}
		if amount == crypto.Amount("0.0") {
			continue
		}
		if amount != crypto.Amount("0.0") {
			found = true
			return found, nil
		}
	}

	return found, nil
}
