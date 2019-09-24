package mixer

import (
	"github.com/Denton24646/gtumbler/pkg/crypto"
	"testing"
)

func TestMixer_CreateDepositAddress(t *testing.T) {
	testMixer := New()
	_, err := testMixer.generateCustomerDepositAddress()
	if err != nil {
		t.Errorf("error generating deposit address: %s", err)
	}
	return
}

func TestMixer_PollDepositAddress(t *testing.T) {
	depositAddress := crypto.Address("Genesis")
	testMixer := New()

	result, err := testMixer.PollDepositAddress(depositAddress)
	if err != nil {
		t.Errorf("error polling deposit address %s", depositAddress)
	}

	t.Logf("address %s contains %s funds", depositAddress, result)
}

func TestMixer_HandleTransaction(t *testing.T) {
	// first generate a new address with some funds in it, simulating a valid deposit
	depositAddr, err := crypto.CreateAddress()
	if err != nil {
		t.Fatalf("error creating address: %s", err)
	}

	// seed newly generated address with funds
	err = crypto.Send(crypto.Address("Genesis"), depositAddr, crypto.Amount("1.0"))
	if err != nil {
		t.Fatalf(" error sending funds: %s", err)
	}

	// create mixer and send funds (after being mixed) back to genesis address
	testMixer := New()
	testMixer.Customers[12] = CustomerData{
		CleanAddresses: []crypto.Address{
			0: "Genesis",
		},
		DepositAddress: depositAddr,
		Fee:            0.05,
	}
	err = testMixer.HandleTransaction(12)
	if err != nil {
		t.Errorf("error handling transaction: %s", err)
	}
}
