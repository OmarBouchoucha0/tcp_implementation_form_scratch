package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
)

func readBytes(packetChan chan []byte, ifce *water.Interface) error {
	packet := make([]byte, 2000)
	n, err := ifce.Read(packet)
	if err != nil {
		return err
	}
	packetChan <- packet[:n]
	return nil
}

type IPHeader struct {
	Version  uint8
	Length   uint8
	TOS      uint8
	TotalLen uint16
	ID       uint16
	Flags    uint8
	FragOff  uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	SrcIP    [4]byte
	DstIP    [4]byte
}

func ipParsing(packet []byte) map[string]interface{} {
	ipHeader := IPHeader{
		Version:  packet[0] >> 4,
		Length:   packet[0] & 0x0f,
		TOS:      packet[1],
		TotalLen: uint16(packet[2])<<8 | uint16(packet[3]),
		ID:       uint16(packet[4])<<8 | uint16(packet[5]),
		Flags:    packet[6] >> 5,
		FragOff:  uint16(packet[6]&0x1f)<<8 | uint16(packet[7]),
		TTL:      packet[8],
		Protocol: packet[9],
		Checksum: uint16(packet[10])<<8 | uint16(packet[11]),
		SrcIP:    [4]byte{packet[12], packet[13], packet[14], packet[15]},
		DstIP:    [4]byte{packet[16], packet[17], packet[18], packet[19]},
	}

	headerMap := map[string]interface{}{
		"Version":  ipHeader.Version,
		"Length":   ipHeader.Length,
		"TOS":      ipHeader.TOS,
		"TotalLen": ipHeader.TotalLen,
		"ID":       ipHeader.ID,
		"Flags":    ipHeader.Flags,
		"FragOff":  ipHeader.FragOff,
		"TTL":      ipHeader.TTL,
		"Protocol": ipHeader.Protocol,
		"Checksum": ipHeader.Checksum,
		"SrcIP":    fmt.Sprintf("%d.%d.%d.%d", ipHeader.SrcIP[0], ipHeader.SrcIP[1], ipHeader.SrcIP[2], ipHeader.SrcIP[3]),
		"DstIP":    fmt.Sprintf("%d.%d.%d.%d", ipHeader.DstIP[0], ipHeader.DstIP[1], ipHeader.DstIP[2], ipHeader.DstIP[3]),
	}

	return headerMap
}

func packetPrint(packet []byte) {
	ipHeader := ipParsing(packet)
	version := ipHeader["Version"].(uint8)
	if version != 4 {
		fmt.Println("wrong version")
		return
	}
	keysInOrder := []string{
		"Version",
		"Length",
		"TOS",
		"TotalLen",
		"ID",
		"Flags",
		"FragOff",
		"TTL",
		"Protocol",
		"Checksum",
		"SrcIP",
		"DstIP",
	}
	fmt.Println("------------------------------------")
	log.Println("New Packet")
	for _, key := range keysInOrder {
		value, ok := ipHeader[key]
		if !ok {
			continue
		}
		fmt.Printf("%s: %v\n", key, value)
	}
	fmt.Println("------------------------------------")
}
