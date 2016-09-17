package chord

import "testing"

const (
	// M3 denotes the use of 3 bits in Chord ID:s.
	M3 = 3
)

func TestFingerStart(t *testing.T) {
	expectNodeFingerStart := func(node *Node, finger int, expected int64) {
		actual := node.Finger(finger).Start()

		if actual.value.Int64() != expected {
			t.Errorf("finger[%d].start %v != expected %v", finger, actual, expected)
		}
	}
	// See figure 3(a), Chord paper.
	{
		a := newNode(stringAddr("a"), newID(1, M3))
		expectNodeFingerStart(a, 1, 2)
		expectNodeFingerStart(a, 2, 3)
		expectNodeFingerStart(a, 3, 5)
	}
}

func TestFingerInterval(t *testing.T) {
	expectNodeFingerInterval := func(node *Node, finger int, expectedStart, expectedStop int64) {
		actualInterval := node.Finger(finger).Interval()
		actualStart := actualInterval.start
		actualStop := actualInterval.stop

		if actualStart.value.Int64() != expectedStart || actualStop.value.Int64() != expectedStop {
			t.Errorf("finger[%d].interval %v != expected [%v,%v)", finger, actualInterval, expectedStart, expectedStop)
		}
	}
	// See figure 3(b), Chord paper.
	{
		a := newNode(stringAddr("a"), newID(0, M3))
		expectNodeFingerInterval(a, 1, 1, 2)
		expectNodeFingerInterval(a, 2, 2, 4)
		expectNodeFingerInterval(a, 3, 4, 0)

		b := newNode(stringAddr("b"), newID(1, M3))
		expectNodeFingerInterval(b, 1, 2, 3)
		expectNodeFingerInterval(b, 2, 3, 5)
		expectNodeFingerInterval(b, 3, 5, 1)

		c := newNode(stringAddr("c"), newID(3, M3))
		expectNodeFingerInterval(c, 1, 4, 5)
		expectNodeFingerInterval(c, 2, 5, 7)
		expectNodeFingerInterval(c, 3, 7, 3)
	}
}
