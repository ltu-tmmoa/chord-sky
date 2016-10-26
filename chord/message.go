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

// Decodes ones message from reader and sets contsnt
func decodeMessage(r io.Reader) (*message, error) {
	m := message{}
	n, err := fmt.Fscanf(r, messageFormat, &m.typ, &m.arg0, &m.arg1)
	if err != nil && n == 0 {
		return nil, err
	}
	return &m, nil
}

func (m *message) encode(w io.Writer) error {
	_, err := fmt.Fprintf(w, messageFormat, m.typ, m.arg0, m.arg1)
	return err
}
