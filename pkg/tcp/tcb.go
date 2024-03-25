package tcp

/*
SND.UNA - send unacknowledged
SND.NXT - send next
SND.WND - send window
SND.UP  - send urgent pointer
SND.WL1 - segment sequence number used for last window update
SND.WL2 - segment acknowledgment number used for last window update
ISS     - initial send sequence number

	     1         2          3          4
	----------|----------|----------|----------
	       SND.UNA    SND.NXT    SND.UNA
	                            +SND.WND

1 - old sequence numbers which have been acknowledged
2 - sequence numbers of unacknowledged data
3 - sequence numbers allowed for new data transmission
4 - future sequence numbers which are not yet allowed

	Send Sequence Space
*/
type SendSequence struct {
	UNA uint32
	NXT uint32
	WND uint16
	UP  bool
	WL1 uint32
	WL2 uint32
	ISS uint32
}

/*
RCV.NXT - receive next
RCV.WND - receive window
RCV.UP  - receive urgent pointer
IRS     - initial receive sequence number


     1          2          3
 ----------|----------|----------
        RCV.NXT    RCV.NXT
                  +RCV.WND

1 - old sequence numbers which have been acknowledged
2 - sequence numbers allowed for new reception
3 - future sequence numbers which are not yet allowed

 Receive Sequence Space
*/

type ReceiveSegment struct {
	NXT uint16
	WND uint32
	UP  bool
	IRS uint32
}

/*
SEG.SEQ - segment sequence number
SEG.ACK - segment acknowledgment number
SEG.LEN - segment length
SEG.WND - segment window
SEG.UP  - segment urgent pointer
SEG.PRC - segment precedence value
*/
type CurrentSegment struct {
	SEQ uint16
	ACK uint16
	LEN uint16
	WND uint16
	UP  uint16
	PRC uint16
}
