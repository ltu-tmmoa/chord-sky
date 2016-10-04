package chord

import (
	"math/big"
	"testing"
)

// stringAddr allows a regular string to be treated as a net.Addr.
type stringAddr string

func (s stringAddr) Network() string {
	return string(s)
}

func (s stringAddr) String() string {
	return string(s)
}

func newHash64(value int64, bits int) *Hash {
	return newHash(*big.NewInt(value), bits)
}

func TestIdIntervalContainsEE(t *testing.T) {
	createIntervalTestFactory := prepareIntervalTester(t, idIntervalContainsEE)
	{
		intervalTestFactory := createIntervalTestFactory(2, 6)
		{
			expectIntervalToContain := intervalTestFactory(true)
			expectIntervalToContain(3)
			expectIntervalToContain(4)
			expectIntervalToContain(5)
		}
		{
			expectIntervalToNotContain := intervalTestFactory(false)
			expectIntervalToNotContain(0)
			expectIntervalToNotContain(2)
			expectIntervalToNotContain(6)
		}
	}
	{
		intervalTestFactory := createIntervalTestFactory(4, 0)
		{
			expectIntervalToContain := intervalTestFactory(true)
			expectIntervalToContain(5)
			expectIntervalToContain(6)
			expectIntervalToContain(7)
		}
		{
			expectIntervalToNotContain := intervalTestFactory(false)
			expectIntervalToNotContain(0)
			expectIntervalToNotContain(2)
			expectIntervalToNotContain(4)
		}
	}
}

func TestIdIntervalContainsEI(t *testing.T) {
	createIntervalTestFactory := prepareIntervalTester(t, idIntervalContainsEI)
	{
		intervalTestFactory := createIntervalTestFactory(2, 6)
		{
			expectIntervalToContain := intervalTestFactory(true)
			expectIntervalToContain(3)
			expectIntervalToContain(4)
			expectIntervalToContain(6)
		}
		{
			expectIntervalToNotContain := intervalTestFactory(false)
			expectIntervalToNotContain(0)
			expectIntervalToNotContain(2)
			expectIntervalToNotContain(7)
		}
	}
	{
		intervalTestFactory := createIntervalTestFactory(4, 0)
		{
			expectIntervalToContain := intervalTestFactory(true)
			expectIntervalToContain(5)
			expectIntervalToContain(6)
			expectIntervalToContain(0)
		}
		{
			expectIntervalToNotContain := intervalTestFactory(false)
			expectIntervalToNotContain(1)
			expectIntervalToNotContain(2)
			expectIntervalToNotContain(4)
		}
	}
}

func TestIdIntervalContainsIE(t *testing.T) {
	createIntervalTestFactory := prepareIntervalTester(t, idIntervalContainsIE)
	{
		intervalTestFactory := createIntervalTestFactory(2, 6)
		{
			expectIntervalToContain := intervalTestFactory(true)
			expectIntervalToContain(2)
			expectIntervalToContain(3)
			expectIntervalToContain(5)
		}
		{
			expectIntervalToNotContain := intervalTestFactory(false)
			expectIntervalToNotContain(0)
			expectIntervalToNotContain(1)
			expectIntervalToNotContain(6)
		}
	}
	{
		intervalTestFactory := createIntervalTestFactory(4, 0)
		{
			expectIntervalToContain := intervalTestFactory(true)
			expectIntervalToContain(4)
			expectIntervalToContain(6)
			expectIntervalToContain(7)
		}
		{
			expectIntervalToNotContain := intervalTestFactory(false)
			expectIntervalToNotContain(0)
			expectIntervalToNotContain(1)
			expectIntervalToNotContain(3)
		}
	}
}

func prepareIntervalTester(t *testing.T, f func(ID, ID, ID) bool) func(start, stop int64) func(bool) func(int64) {
	return func(start, stop int64) func(bool) func(int64) {
		return func(contains bool) func(int64) {
			return func(other int64) {
				start0 := newHash64(start, M3)
				stop0 := newHash64(stop, M3)
				other0 := newHash64(other, M3)
				if f(start0, stop0, other0) != contains {
					var operator string
					if contains {
						operator = "contains"
					} else {
						operator = "not contains"
					}
					t.Errorf("[%v, %v) %v %v", start, stop, operator, other)
				}
			}
		}
	}
}
