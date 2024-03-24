package main

import (
	"errors"
	"fmt"
	"github.com/songgao/water"
	"log"
)

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
	Options  []byte
	Padding  []byte
}

func readBytes(packetChan chan []byte, ifce *water.Interface) error {
	packet := make([]byte, 2000)
	n, err := ifce.Read(packet)
	if err != nil {
		return err
	}
	packetChan <- packet[:n]
	return nil
}

func validIpPacket(packet []byte) error {
	if len(packet) < 20 {
		return errors.New("not a valid Ip Packet")
	}
	return nil
}

func parsingIpPacket(packet []byte) (IPHeader, error) {
	err := validIpPacket(packet)
	if err != nil {
		return IPHeader{}, err
	}
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
	/*
		if uint8(ipHeader.Length) > 5 {
			ipHeader.Options = packet[20 : uint8(packet[0]&0x0f)*8]
			ipHeader.Padding = packet[uint8(packet[0]&0x0f)*8 : uint16(packet[2])<<8|uint16(packet[3])]
		}
	*/
	return ipHeader, nil
}

func printIpPacket(packet []byte) {
	ipHeader, err := parsingIpPacket(packet)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("New Packet")
	fmt.Printf("Version: %d\n", ipHeader.Version)
	fmt.Printf("Header Length: %d\n", ipHeader.Length)
	fmt.Printf("Type of Service (TOS): %d\n", ipHeader.TOS)
	fmt.Printf("Total Length: %d\n", ipHeader.TotalLen)
	fmt.Printf("Identification: %d\n", ipHeader.ID)
	fmt.Printf("Flags: %d\n", ipHeader.Flags)
	fmt.Printf("Fragment Offset: %d\n", ipHeader.FragOff)
	fmt.Printf("Time to Live (TTL): %d\n", ipHeader.TTL)
	fmt.Printf("Protocol: %d\n", ipHeader.Protocol)
	fmt.Printf("Checksum: %d\n", ipHeader.Checksum)
	fmt.Printf("Source IP: %d.%d.%d.%d\n", ipHeader.SrcIP[0], ipHeader.SrcIP[1], ipHeader.SrcIP[2], ipHeader.SrcIP[3])
	fmt.Printf("Destination IP: %d.%d.%d.%d\n", ipHeader.DstIP[0], ipHeader.DstIP[1], ipHeader.DstIP[2], ipHeader.DstIP[3])
	fmt.Println("------------------------------------")
}
