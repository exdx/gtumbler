package mixer

import (
	"encoding/json"
	"github.com/Denton24646/gtumbler/pkg/crypto"
	"github.com/Denton24646/gtumbler/pkg/mixer/tumbler"
	"github.com/Denton24646/gtumbler/pkg/models"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// The mixer is an http server responsible for mixing the client coins by doing the following
// 1. On startup, preseed a certain amount of addresses with coins (to bootstrap the mixing process)
// Since there is no API method for creating coins from scratch, these addresses need to be made from the UI
// The mixer will assume these addresses already exist and are funded
// 2. Accept requests from clients that conform to certain rules (min, max, fee, etc)
// 3. Provide a deposit address back to the client
// 4. Check the blockchain to see if/when the client sends funds to the deposit address
// 5. When funds are received, move funds into smaller random amounts into addresses mixer controls
// 6. Send those funds back to the clients specified address in random intervals

type Server interface {
	// Create is the /create endpoint for the mixer - it accepts the list of new addresses and returns the deposit address
	Create(w http.ResponseWriter, req *http.Request)
	//CreateDepositAddress generates a new deposit address for the customer
	generateCustomerDepositAddress() (crypto.Address, error)
	//PollDepositAddress checks the deposit address periodically to see if the client deposited funds
	PollDepositAddress(address crypto.Address) (crypto.Amount, error)
	// HandleTransaction is responsible for all the backend work of the mixer service
	HandleTransaction(id int) error
}

type Mixer struct {
	// Customer Ids is an map of Ids of clients to their ultimate clean addresses and deposit information
	// It's an in memory datastore for the purposes of having access to customer information
	// TODO these ids would be used to further obfuscate in the mixing process
	// TODO use more durable datastore
	Customers map[int]CustomerData
	// house addresses is an array of addresses the house owns and are already funded
	// these addresses can be used by the tumbler, which has no knowledge of the mixer and simply moves coins around
	HouseAddresses []crypto.Address
}

type CustomerData struct {
	CleanAddresses  []crypto.Address
	DepositAddress  crypto.Address
	Fee             float64
}

func New() *Mixer {
	return &Mixer{
		HouseAddresses: []crypto.Address{
			0: "House1",
			1: "House2",
			2: "House3",
			3: "House4",
			4: "House5",
		},
	}
}

func (m *Mixer) Create(w http.ResponseWriter, req *http.Request) {
	request := &models.CleanAddressRequest{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(404)
	}
	defer req.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(404)
	}

	depositAddress, err := m.generateCustomerDepositAddress()
	if err != nil {
		w.WriteHeader(404)
	}

	customerId := request.Id
	m.Customers[customerId] = CustomerData{
		CleanAddresses: request.Addresses,
		DepositAddress: depositAddress,
		Fee: rand.Float64() * 0.01,
	}

	response := &models.CleanAddressResponse{
		DepositAddress: depositAddress,
	}
	res,_ := json.Marshal(response)

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(404)
	}

	// concurrently handle customer transactions via goroutines
	go func(id int) {
		err := m.HandleTransaction(id)
		if err != nil {
			w.WriteHeader(404)
		}
	}(customerId)

	w.WriteHeader(http.StatusOK)
}

// generateCustomerDepositAddress generates new addresses for customers to deposit into
func (m *Mixer) generateCustomerDepositAddress() (crypto.Address, error) {
	address, err := crypto.CreateAddress()
	if err != nil {
		return "", err
	}
	return address, nil
}

// PollDepositAddress checks the provided addresses to see if the customer deposited funds yet
func (m *Mixer) PollDepositAddress(address crypto.Address) (crypto.Amount, error) {
	amount, err := crypto.CheckAddress(address)
	if err != nil {
		return "0", err
	}
	if amount == crypto.Amount("0.00") {
		return "0", nil
	}
	return amount, nil
}

// HandleTransaction is the controller that handles the flow of customer funds
// First it polls to check the customer deposit address for funds
// Once funds are sent it uses the tumbler to tumble funds and send them back to the mixer
func (m *Mixer) HandleTransaction(id int) error {
	var amount crypto.Amount
	for {
		deposit, err := m.PollDepositAddress(m.Customers[id].DepositAddress)
		if err != nil {
			return err
		}
		if deposit != crypto.Amount("0") {
			amount = deposit
			break
		}
		time.Sleep(10 * time.Second)
	}

	log.Printf(" **** Received %s coins from address %s with return addresses %v", amount, m.Customers[id].DepositAddress,
		m.Customers[id].CleanAddresses)

	tumblr := tumbler.New(amount)
	err := tumblr.Mix(m.Customers[id].DepositAddress, m.HouseAddresses)
	if err != nil {
		return err
	}

	log.Printf("**** Tumbled coins from %s to house addresses %s successfully", m.Customers[id].DepositAddress,
		m.HouseAddresses)

	err = tumblr.SendMixedFunds(m.Customers[id].CleanAddresses, m.HouseAddresses)
	if err != nil {
		return err
	}

	log.Printf("**** Sent mixed coins back to %s successfully ****", m.Customers[id].CleanAddresses)

	return nil
}

