package main

import (
	"github.com/Denton24646/gtumbler/pkg/mixer"
	"log"
	"net/http"
)

func main() {
	log.Print("**** Starting gtumbler mixer service ****")
	m := mixer.New()

	http.HandleFunc("/create", m.Create)
	log.Print("**** Listening on port 8989 for new mixer deposit transactions ****")
	log.Fatal(http.ListenAndServe(":8989", nil))
}
