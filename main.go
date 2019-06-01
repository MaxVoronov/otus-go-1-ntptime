package main

import (
	"log"
	"time"
)

const server = "time.apple.com"
const port = "123"

func main() {
	ntp := &Ntp{server, port, 5}

	log.Printf("Local time: %s", time.Now())
	log.Printf("NTP Time: %s", ntp.GetTime())
}
