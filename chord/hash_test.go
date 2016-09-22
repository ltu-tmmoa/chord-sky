package chord

import "testing"

func TestHashCmp(t *testing.T) {
	a := newHash64(5, M3)
	b := newHash64(4, M3)

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
	a := newHash64(5, M3)
	b := newHash64(1, M3)

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
