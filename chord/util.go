package chord

func idIntervalContainsEE(start, stop, other ID) bool {
	a := other.Cmp(start)
	b := other.Cmp(stop)

	if start.Cmp(stop) < 0 {
		return a > 0 && b < 0
	}
	return a > 0 || b < 0
}

func idIntervalContainsEI(start, stop, other ID) bool {
	a := other.Cmp(start)
	b := other.Cmp(stop)

	if start.Cmp(stop) < 0 {
		return a > 0 && b <= 0
	}
	return a > 0 || b <= 0
}

func idIntervalContainsIE(start, stop, other ID) bool {
	a := other.Cmp(start)
	b := other.Cmp(stop)

	if start.Cmp(stop) < 0 {
		return a >= 0 && b < 0
	}
	return a >= 0 || b < 0
}
