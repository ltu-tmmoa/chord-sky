package chord

import (
	"fmt"
	"io"
)

const (
	messageTypeHeartbeat messageType = iota

	messageTypeGetFingerNode
	messageTypeGetSuccessor
	//messageTypeGetSuccessorList
	messageTypeGetPredecessor

	messageTypeFindSuccessor
	messageTypeFindPredecessor

	messageTypeSetFingerNode
	messageTypeSetSuccessor
	messageTypeSetPredecessor

	//messageTypeLeave
)

const (
	messageFormat = "CHORD %d %s %s\r\n"
)

// Indicates the intent of some message.
type messageType int

// Some message.
type message struct {
	typ  messageType
	arg0 string
	arg1 string
}

// Reads and decodes one message from reader.
func decodeMessage(r io.Reader, typ messageType) (*message, error) {
	m := message{}
	n, err := fmt.Fscanf(r, messageFormat, &m.typ, &m.arg0, &m.arg1)
	if err != nil && n == 0 {
		return nil, err
	}
	if m.typ != typ {
		return nil, fmt.Errorf("Decoded message of type %d, expected %d.", m.typ, typ)
	}
	return &m, nil
}

// Encodes and writes message to writer.
func encodeMessage(w io.Writer, typ messageType, arg0, arg1 string) error {
	m := message{
		typ:  typ,
		arg0: arg0,
		arg1: arg1,
	}
	return m.encode(w)
}

// Encodes and writes this message to writer.
func (m *message) encode(w io.Writer) error {
	_, err := fmt.Fprintf(w, messageFormat, m.typ, m.arg0, m.arg1)
	return err
}
