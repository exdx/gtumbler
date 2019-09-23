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
	_, err = c.CreateCleanAddresses(config.NumberAddresses)
	if err != nil {
		log.Printf("error generating addresses: %s", err )
	}
	fmt.Println(" **** Address creation successful ****")

	err = c.SendCleanAddresses()
	if err != nil {
		log.Printf("\nerror receiving deposit address: %s", err)
	}
	fmt.Printf("**** gtumbler deposit address %s\n", c.DepositAddress)
	fmt.Printf("**** Sending amount %s to deposit address from address %s\n", config.Size, config.SendAddress)

	err = c.SendDeposit(config.SendAddress, config.Size)
	if err != nil {
		log.Printf(" error sending funds to deposit address %s: %s", c.DepositAddress, err)
	}
	fmt.Println("**** Deposit sent to gtumbler mixer ****")
	fmt.Println("**** Waiting 5 seconds and checking blockchain for mixing status ****")

	for {
		finished, err := c.CheckCleanAddresses()
		if err != nil {
			fmt.Println("**** Issue mixing your coins. Check the blockchain for more information. ****")
		}
		if finished {
			fmt.Println("**** Successful mixing. The coins are now in the addresses specified. Thank you for using gtumbler. ****")
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
