package data

import (
	"math/big"
	"testing"
)

const (
	bits = 3
)

func newID64(value int64, bits int) *ID {
	return NewID(big.NewInt(value), bits)
}

func TestHashCmp(t *testing.T) {
	a := newID64(5, bits)
	b := newID64(4, bits)

	if r := a.Cmp(a); r != 0 {
		t.Errorf("%v.Cmp(%v) %v != 0", a, a, r)
	}
	if r := a.Cmp(b); r != 1 {
		t.Errorf("%v.Cmp(%v) %v != 1", a, b, r)
	}
	if r := b.Cmp(b); r != 0 {
		t.Errorf("%v.Cmp(%v) %v != 0", b, b, r)
	}
	if r := b.Cmp(a); r != -1 {
		t.Errorf("%v.Cmp(%v) %v != -1", b, a, r)
	}
}

func TestHashDiff(t *testing.T) {
	a := newID64(5, bits)
	b := newID64(1, bits)

	if r := a.Diff(a).BigInt().Int64(); r != 0 {
		t.Errorf("%v.Diff(%v) %v != 0", a, a, r)
	}
	if r := a.Diff(b).BigInt().Int64(); r != 4 {
		t.Errorf("%v.Diff(%v) %v != 3", a, b, r)
	}
	if r := b.Diff(b).BigInt().Int64(); r != 0 {
		t.Errorf("%v.Diff(%v) %v != 0", b, b, r)
	}
	if r := b.Diff(a).BigInt().Int64(); r != 4 {
		t.Errorf("%v.Diff(%v) %v != 4", b, a, r)
	}
}

func TestParseID(t *testing.T) {
	testParseID3 := func(t *testing.T, s string, expectedID int64) {
		id, ok := ParseID(s, bits)
		if !ok {
			t.Errorf("Failed to parse ID: %s", s)
			return
		}
		if id.BigInt().Int64() != expectedID {
			t.Errorf("ID %v != %d", id, expectedID)
		}
	}

	testParseID3(t, "-1", 7)
	testParseID3(t, "0", 0)
	testParseID3(t, "2", 2)
	testParseID3(t, "7", 7)
	testParseID3(t, "9", 1)
}
