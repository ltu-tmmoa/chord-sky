package chord

import "fmt"

func verifyIndexOrPanic(len, i int) {
	if 1 > i || i > len {
		panic(fmt.Sprintf("%d not in [1,%d]", i, len))
	}
}
