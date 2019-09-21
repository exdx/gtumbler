package main

import (
	"fmt"
	"github.com/Denton24646/gtumbler/pkg/client"
	"github.com/crgimenes/goconfig"
	"log"
	"time"
)

func main() {
	// get configuration from the command line or the environment
	config := client.Config{}
	err := goconfig.Parse(&config)
	if err != nil {
		log.Fatalf("parsing config: %s", err)
	}

	// setup user client
	fmt.Println("**** Welcome to the gtumber client ****")
	c := client.New(config)

	// create addresses or use addresses provided
	fmt.Println("**** Generating newly created address for use with the gtumbler mixer")
	c.CreateCleanAddresses(config.NumberAddresses)

	// show deposit address
	fmt.Println(" **** Address creation successful ****")
	fmt.Printf("**** gtumbler deposit address %s", c.DepositAddress)
	fmt.Println("\n **** Sending amount %s to deposit address from address %s", config.Size, config.SendAddress)
	c.SendDeposit(config.SendAddress, config.Size)

	fmt.Println("**** Deposit sent to gtumbler mixer ****")
	fmt.Println("**** Waiting 5 seconds and checking blockchain for mixing status ****")

	for {
		finished, err := c.CheckCleanAddresses()
		if err != nil {
			fmt.Println("**** Issue mixing your coins. Check the blockchain for more information. ****")
		}
		if finished {
			fmt.Println("**** Successful mixing. The coins are now in the addresses specified ****")
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
