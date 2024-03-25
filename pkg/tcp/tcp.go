package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/OmarBouchoucha0/tcp_implementation_from_scratch/pkg/ip"
	"log"
)

type TCPHeader struct {
	/*
	   0                   1                   2                   3
	    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |          Source Port          |       Destination Port        |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                        Sequence Number                        |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                    Acknowledgment Number                      |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |  Data |           |U|A|P|R|S|F|                               |
	   | Offset| Reserved  |R|C|S|S|Y|I|            Window             |
	   |       |           |G|K|H|T|N|N|                               |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |           Checksum            |         Urgent Pointer        |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                    Options                    |    Padding    |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                             data                              |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
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
	data                 []byte
}

type Flags struct {
	/*
	   URG:  Urgent Pointer field significant
	   ACK:  Acknowledgment field significant
	   PSH:  Push Function
	   RST:  Reset the connection
	   SYN:  Synchronize sequence numbers
	   FIN:  No more data from sender
	*/
	FIN bool
	SYN bool
	RST bool
	PSH bool
	ACK bool
	URG bool
}

func IsTcpPacket(ipHeader ip.IPHeader) bool {
	version := ipHeader.Version
	protocol := ipHeader.Protocol
	if version != 4 && protocol != 6 {
		return false
	}
	return true
}
func UnparseTcpPacket(tcpHeader TCPHeader) ([]byte, error) {
	totalLength := 20 + len(tcpHeader.data)
	packet := make([]byte, totalLength)
	binary.BigEndian.PutUint16(packet[0:2], tcpHeader.SourcePort)
	binary.BigEndian.PutUint16(packet[2:4], tcpHeader.DestinationPort)
	binary.BigEndian.PutUint32(packet[4:8], tcpHeader.SequenceNumber)
	binary.BigEndian.PutUint32(packet[8:12], tcpHeader.AcknowledgmentNumber)
	packet[12] = tcpHeader.DataOffset << 4
	packet[12] |= tcpHeader.Reserved
	flags := uint16(0)
	if tcpHeader.Flags.FIN {
		flags |= 1 << 0
	}
	if tcpHeader.Flags.SYN {
		flags |= 1 << 1
	}
	if tcpHeader.Flags.RST {
		flags |= 1 << 2
	}
	if tcpHeader.Flags.PSH {
		flags |= 1 << 3
	}
	if tcpHeader.Flags.ACK {
		flags |= 1 << 4
	}
	if tcpHeader.Flags.URG {
		flags |= 1 << 5
	}
	binary.BigEndian.PutUint16(packet[12:14], flags)
	binary.BigEndian.PutUint16(packet[14:16], tcpHeader.WindowSize)
	binary.BigEndian.PutUint16(packet[16:18], tcpHeader.Checksum)
	binary.BigEndian.PutUint16(packet[18:20], tcpHeader.UrgentPointer)

	copy(packet[20:], tcpHeader.data)

	return packet, nil
}

func ParsingTcpPacket(packet []byte) (TCPHeader, error) {
	ipHeader, err := ip.ParsingIpPacket(packet)
	if err != nil {
		return TCPHeader{}, err
	}
	if IsTcpPacket(ipHeader) {
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

func PrintTcpPacket(packet []byte) {
	tcpHeader, err := ParsingTcpPacket(packet)
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
