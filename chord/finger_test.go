package chord

import "testing"

const (
	// M3 denotes the use of 3 bits in Chord ID:s.
	M3 = 3
)

func TestFingerStart(t *testing.T) {
	expectNodeFingerStart := func(node *LocalNode, finger int, expected int64) {
		actual := node.Finger(finger).Start()

		if actual.BigInt().Int64() != expected {
			t.Errorf("finger[%d].start %v != expected %v", finger, actual, expected)
		}
	}
	// See figure 3(a), Chord paper.
	{
		a := newLocalNode(fakeAddr(1), newHash64(1, M3))
		expectNodeFingerStart(a, 1, 2)
		expectNodeFingerStart(a, 2, 3)
		expectNodeFingerStart(a, 3, 5)
	}
}

func TestFingerInterval(t *testing.T) {
	expectNodeFingerInterval := func(node *LocalNode, finger int, expectedStart, expectedStop int64) {
		actualInterval := node.Finger(finger).Interval()

		if actualInterval.Start().BigInt().Int64() != expectedStart || actualInterval.Stop().BigInt().Int64() != expectedStop {
			t.Errorf("finger[%d].interval [%v,%v) != expected [%v,%v)", finger, actualInterval.Start(), actualInterval.Stop(), expectedStart, expectedStop)
		}
	}
	// See figure 3(b), Chord paper.
	{
		a := newLocalNode(fakeAddr(0), newHash64(0, M3))
		expectNodeFingerInterval(a, 1, 1, 2)
		expectNodeFingerInterval(a, 2, 2, 4)
		expectNodeFingerInterval(a, 3, 4, 0)

		b := newLocalNode(fakeAddr(1), newHash64(1, M3))
		expectNodeFingerInterval(b, 1, 2, 3)
		expectNodeFingerInterval(b, 2, 3, 5)
		expectNodeFingerInterval(b, 3, 5, 1)

		c := newLocalNode(fakeAddr(3), newHash64(3, M3))
		expectNodeFingerInterval(c, 1, 4, 5)
		expectNodeFingerInterval(c, 2, 5, 7)
		expectNodeFingerInterval(c, 3, 7, 3)
	}
}
