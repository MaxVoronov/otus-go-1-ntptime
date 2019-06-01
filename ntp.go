package main

import (
	"encoding/binary"
	"log"
	"net"
	"time"
)

const ntpEpochOffset uint64 = 2208988800
const ntpDefaultTimeout int = 5

type Ntp struct {
	Server  string
	Port    string
	Timeout time.Duration
}

// NTP packet format (v3 with optional v4 fields removed)
// Protocol details: https://ru.wikipedia.org/wiki/NTP
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |LI | VN  |Mode |    Stratum     |     Poll      |  Precision   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         Root Delay                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         Root Dispersion                       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Reference ID                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                     Reference Timestamp (64)                  +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Origin Timestamp (64)                    +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Receive Timestamp (64)                   +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Transmit Timestamp (64)                  +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type NtpPackage struct {
	Settings       uint8 // Leap Indicator [2 bits] + Version Number [3 bits] + Mode [3 bits]
	Stratum        uint8
	Poll           uint8
	Precision      uint8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceId    uint32
	ReferenceSec   uint32
	ReferenceNano  uint32
	OriginSec      uint32
	OriginNano     uint32
	ReceiveSec     uint32
	ReceiveNano    uint32
	TransmitSec    uint32
	TransmitNano   uint32
}

func (ntp *Ntp) GetTime() time.Time {
	ntpTime := ntp.sendRequest()
	seconds := int64(uint64(ntpTime.TransmitSec) - ntpEpochOffset)
	nano := int64(ntpTime.TransmitNano) * 1e9 >> 32

	return time.Unix(seconds, nano)
}

func (ntp *Ntp) sendRequest() *NtpPackage {
	client, err := net.Dial("udp", net.JoinHostPort(ntp.Server, ntp.Port))
	if err != nil {
		log.Fatalf("Connection error: %s", err)
	}
	defer client.Close()

	timeout := time.Duration(ntp.Timeout)
	if timeout == 0 {
		timeout = time.Duration(ntpDefaultTimeout)
	}

	if err = client.SetDeadline(time.Now().Add(timeout * time.Second)); err != nil {
		log.Fatalf("Failed to set deadline: %s", err)
	}

	request := &NtpPackage{Settings: 0x1B}
	if err = binary.Write(client, binary.BigEndian, request); err != nil {
		log.Fatalf("Failed to send request: %s", err)
	}

	response := &NtpPackage{}
	if err = binary.Read(client, binary.BigEndian, response); err != nil {
		log.Fatalf("Failed to read response: %s", err)
	}

	return response
}
