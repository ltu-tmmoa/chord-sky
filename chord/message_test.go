package chord

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestMessageTranscoding(t *testing.T) {
	var buf bytes.Buffer

	checkErr := func(err error) {
		if err != nil {
			panic(fmt.Sprintf("Unexpected error: %v", err))
		}
	}
	encode := func(typ messageType, arg0, arg1 string) {
		m := message{
			typ:  typ,
			arg0: arg0,
			arg1: arg1,
		}
		err := m.encode(&buf)
		checkErr(err)
	}
	decode := func(typ int, arg0, arg1 string) {
		m, err := decodeMessage(&buf)
		checkErr(err)
		if m.typ != messageType(typ) {
			t.Errorf("Message typ is %d, expected %d.", m.typ, typ)
		}
		if m.arg0 != arg0 {
			t.Errorf("Message arg0 is \"%s\", expected \"%s\".", m.arg0, arg0)
		}
		if m.arg1 != arg1 {
			t.Errorf("Message arg1 is \"%s\", expected \"%s\".", m.arg1, arg1)
		}
	}

	encode(0, "", "")
	encode(1337, "monkey", "banana")
	encode(7331, "horsie", "applie")
	encode(-129, "power", "")
	encode(1, "", "")

	decode(0, "", "")
	decode(1337, "monkey", "banana")
	decode(7331, "horsie", "applie")
	decode(-129, "power", "")
	decode(1, "", "")

	{
		if _, err := decodeMessage(&buf); err != io.ErrUnexpectedEOF {
			t.Errorf("Expected ErrUnexpectedEOF, got %v.", err)
		}
	}
}
