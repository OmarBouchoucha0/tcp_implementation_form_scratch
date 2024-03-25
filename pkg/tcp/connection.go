package tcp

import (
	"errors"
)

/*
   LISTEN - represents waiting for a connection request from any remote
   TCP and port.

   SYN-SENT - represents waiting for a matching connection request
   after having sent a connection request.

   SYN-RECEIVED - represents waiting for a confirming connection
   request acknowledgment after having both received and sent a
   connection request.

   ESTABLISHED - represents an open connection, data received can be
   delivered to the user.  The normal state for the data transfer phase
   of the connection.

   FIN-WAIT-1 - represents waiting for a connection termination request
   from the remote TCP, or an acknowledgment of the connection
   termination request previously sent.

   FIN-WAIT-2 - represents waiting for a connection termination request
   from the remote TCP.

   CLOSE-WAIT - represents waiting for a connection termination request
   from the local user.

   CLOSING - represents waiting for a connection termination request
   acknowledgment from the remote TCP.

   LAST-ACK - represents waiting for an acknowledgment of the
   connection termination request previously sent to the remote TCP
   (which includes an acknowledgment of its connection termination
   request).

   TIME-WAIT - represents waiting for enough time to pass to be sure
   the remote TCP received the acknowledgment of its connection
   termination request.

   CLOSED - represents no connection state at all.
*/

type Connection struct {
	LISTEN       bool
	SYN_SENT     bool
	SYN_RECEIVED bool
	ESTABLISHED  bool
	FIN_WAIT_1   bool
	FIN_WAIT_2   bool
	CLOSE_WAIT   bool
	CLOSING      bool
	LAST_ACK     bool
	CLOSED       bool
}

func acknowledgment(tcpHeader TCPHeader) errors {
	recivedSegment := ReceiveSegment{
		IRS: tcpHeader.SequenceNumber,
		WND: tcpHeader.SequenceNumber + 1,
		NXT: tcpHeader.WindowSize,
		UP:  false,
	}
	sentSengment := SendSequence{
		ISS: 0,
		WND: 10,
	}
	sentSengment.UP = false
	sentSengment.UNA = sentSengment.ISS
	sentSengment.NXT = sentSengment.UNA + 1

}
