package chord

import (
	"fmt"

	"github.com/ltu-tmmoa/chord-sky/data"
)

func idIntervalContainsEE(start, stop, other *data.ID) bool {
	a := other.Cmp(start)
	b := other.Cmp(stop)

	if start.Cmp(stop) < 0 {
		return a > 0 && b < 0
	}
	return a > 0 || b < 0
}

func idIntervalContainsEI(start, stop, other *data.ID) bool {
	a := other.Cmp(start)
	b := other.Cmp(stop)

	if start.Cmp(stop) < 0 {
		return a > 0 && b <= 0
	}
	return a > 0 || b <= 0
}

func idIntervalContainsIE(start, stop, other *data.ID) bool {
	a := other.Cmp(start)
	b := other.Cmp(stop)

	if start.Cmp(stop) < 0 {
		return a >= 0 && b < 0
	}
	return a >= 0 || b < 0
}

func verifyIndexOrPanic(len, i int) {
	if 1 > i || i > len {
		panic(fmt.Sprintf("%d not in [1,%d]", i, len))
	}
}
