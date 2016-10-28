package chord

import "testing"

func TestHashCmp(t *testing.T) {
	a := newID64(5, M3)
	b := newID64(4, M3)

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
	a := newID64(5, M3)
	b := newID64(1, M3)

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
		id, ok := parseID(s, M3)
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
