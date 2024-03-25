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

type SendSequence struct {
	/*
	   Send Sequence Variables

	     SND.UNA - send unacknowledged
	     SND.NXT - send next
	     SND.WND - send window
	     SND.UP  - send urgent pointer
	     SND.WL1 - segment sequence number used for last window update
	     SND.WL2 - segment acknowledgment number used for last window
	               update
	     ISS     - initial send sequence number
	*/
	UNA uint16
	NXT uint16
	WND uint16
	UP  bool
	WL1 uint16
	WL2 uint16
	ISS uint16
}

type ReceiveSegment struct {
	/*
	   RCV.NXT - receive next
	   RCV.WND - receive window
	   RCV.UP  - receive urgent pointer
	   IRS     - initial receive sequence number
	*/
	NXT uint16
	WND uint16
	UP  bool
	IRS uint16
}

type CurrentSegment struct {
	/*
	   SEG.SEQ - segment sequence number
	   SEG.ACK - segment acknowledgment number
	   SEG.LEN - segment length
	   SEG.WND - segment window
	   SEG.UP  - segment urgent pointer
	   SEG.PRC - segment precedence value
	*/
	SEQ uint16
	ACK uint16
	LEN uint16
	WND uint16
	UP  uint16
	PRC uint16
}

func IsTcpPacket(ipHeader ip.IPHeader) bool {
	version := ipHeader.Version
	protocol := ipHeader.Protocol
	if version != 4 && protocol != 6 {
		return false
	}
	return true
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
