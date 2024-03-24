package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

type TCPHeader struct {
	SourcePort           uint16
	DestinationPort      uint16
	SequenceNumber       uint32
	AcknowledgmentNumber uint32
	DataOffset           uint8
	Reserved             uint8
	Flags                Flags
	WindowSize           uint16
	Checksum             uint16
	UrgentPointer        uint16
}

type Flags struct {
	FIN bool
	SYN bool
	RST bool
	PSH bool
	ACK bool
	URG bool
}

func isTcpPacket(ipHeader IPHeader) bool {
	version := ipHeader.Version
	protocol := ipHeader.Protocol
	if version != 4 && protocol != 6 {
		return false
	}
	return true
}

func parsingTcpPacket(packet []byte) (TCPHeader, error) {
	ipHeader, err := parsingIpPacket(packet)
	if err != nil {
		return TCPHeader{}, err
	}
	if isTcpPacket(ipHeader) {
		ipHeaderLen := uint16(ipHeader.Length) * 4
		tcpPacket := packet[ipHeaderLen:]
		tcpHeader := TCPHeader{
			SourcePort:           binary.BigEndian.Uint16(tcpPacket[0:2]),
			DestinationPort:      binary.BigEndian.Uint16(tcpPacket[2:4]),
			SequenceNumber:       binary.BigEndian.Uint32(tcpPacket[4:8]),
			AcknowledgmentNumber: binary.BigEndian.Uint32(tcpPacket[8:12]),
			DataOffset:           tcpPacket[12] >> 4,
			Reserved:             tcpPacket[12] & 0x0F,
			WindowSize:           binary.BigEndian.Uint16(tcpPacket[14:16]),
			Checksum:             binary.BigEndian.Uint16(tcpPacket[16:18]),
			UrgentPointer:        binary.BigEndian.Uint16(tcpPacket[18:20]),
			Flags: Flags{
				FIN: tcpPacket[14]&1 != 0,
				SYN: tcpPacket[13]&2 != 0,
				RST: tcpPacket[13]&4 != 0,
				PSH: tcpPacket[13]&8 != 0,
				ACK: tcpPacket[13]&16 != 0,
				URG: tcpPacket[13]&32 != 0,
			},
		}
		return tcpHeader, nil
	} else {
		return TCPHeader{}, errors.New("not a TCP packet")
	}
}

func printTcpPacket(packet []byte) {
	tcpHeader, err := parsingTcpPacket(packet)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("TCP Packet:")
	fmt.Printf("Source Port: %d\n", tcpHeader.SourcePort)
	fmt.Printf("Destination Port: %d\n", tcpHeader.DestinationPort)
	fmt.Printf("Sequence Number: %d\n", tcpHeader.SequenceNumber)
	fmt.Printf("Acknowledgment Number: %d\n", tcpHeader.AcknowledgmentNumber)
	fmt.Printf("Data Offset: %d\n", tcpHeader.DataOffset)
	fmt.Printf("Reserved: %d\n", tcpHeader.Reserved)
	fmt.Printf("Flags: %+v\n", tcpHeader.Flags)
	fmt.Printf("Window Size: %d\n", tcpHeader.WindowSize)
	fmt.Printf("Checksum: %d\n", tcpHeader.Checksum)
	fmt.Printf("Urgent Pointer: %d\n", tcpHeader.UrgentPointer)
	fmt.Println("---------------------------------------------------------")
}
