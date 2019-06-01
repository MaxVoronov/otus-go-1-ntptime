package main

import (
	"log"
	"time"
)

const server = "time.apple.com"
const port = "123"

func main() {
	ntp := NewNtp(server, port, 5)

	ntpTime, err := ntp.GetTime()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	log.Printf("Local time: %s", time.Now())
	log.Printf("NTP Time: %s", ntpTime)
}
