package main

import (
	"log"
	"net/http"

	"github.com/imdevinc/steelseries-stream-toggle/internals/steelseries"
)

var httpClient = &http.Client{}

func main() {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig.InsecureSkipVerify = true
	httpClient.Transport = customTransport
	ss := steelseries.New(httpClient)
	newMode, err := ss.ToggleStreamerMode()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(newMode)
}
