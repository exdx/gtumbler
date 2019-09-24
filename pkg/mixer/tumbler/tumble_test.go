package tumbler

import (
	"fmt"
	"github.com/Denton24646/gtumbler/pkg/crypto"
	"strconv"
	"testing"
)

func TestTumbler_ValidDeposit(t *testing.T) {
	tableTests := []struct{
		amount float64
		valid bool
	}{
		{1.0, true},
		{0.5, true},
		{8, true},
		{17, false},
		{0, false},
		{100, false},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := valid(tt.amount)
			if s != tt.valid {
				t.Errorf("record %d got %t, want %t", i, s, tt.valid)
			}
		})
	}
}

func TestTumbler_Mix(t *testing.T) {
	deposit := crypto.Amount("1.0")
	depositAddr := crypto.Address("Genesis")
	houseAddr := []crypto.Address{
		0: "House1",
		1: "House2",
		2: "House3",
		3: "House4",
		4: "House5",
	}
	testTumbler := New(deposit)
	// check balance of deposit address before test
	// we expect the balance to be lower by one after the test
	// note: this had to be relaxed to be less than some amount because tests run concurrently
	balance, err := crypto.CheckAddress(depositAddr)
	if err != nil {
		t.Errorf("error checking balance of address %s: %s", depositAddr, err)
	}
	err = testTumbler.Mix(depositAddr, houseAddr)
	if err != nil {
		t.Errorf("error mixing coins: %s", err)
	}

	newBalance, err := crypto.CheckAddress(depositAddr)
	if err != nil {
		t.Errorf("error checking balance of address %s: %s", depositAddr, err)
	}

	b, _ :=  strconv.ParseFloat(string(balance), 64)
	nb, _ := strconv.ParseFloat(string(newBalance), 64)
	diff := nb - b

	if diff > 0 {
		t.Errorf("expected difference in deposit address of atleast %f, got %f", 1.0, diff)
	}
}

func TestTumbler_SendMixedFunds(t *testing.T) {
	deposit := crypto.Amount("1.0")
	fundAddr := []crypto.Address{
		0: "Genesis",
	}

	houseAddr := []crypto.Address{
		0: "House1",
		1: "House2",
		2: "House3",
		3: "House4",
		4: "House5",
	}
	testTumbler := New(deposit)
	// check balance of customer address before test
	// we expect the balance to be higher by one after the test
	// note: this had to be relaxed to be greater than some amount because tests run concurrently
	balance, err := crypto.CheckAddress(fundAddr[0])
	if err != nil {
		t.Errorf("error checking balance of address %s: %s", fundAddr, err)
	}
	err = testTumbler.SendMixedFunds(fundAddr, houseAddr)
	if err != nil {
		t.Errorf("error mixing coins: %s", err)
	}

	newBalance, err := crypto.CheckAddress(fundAddr[0])
	if err != nil {
		t.Errorf("error checking balance of address %s: %s", fundAddr, err)
	}

	b, _ :=  strconv.ParseFloat(string(balance), 64)
	nb, _ := strconv.ParseFloat(string(newBalance), 64)
	diff := nb - b

	if diff < 0 {
		t.Errorf("expected difference in deposit address of atleast %f, got %f", 1.0, diff)
	}
}
