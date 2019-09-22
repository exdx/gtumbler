package tumbler

import "math/rand"

// strategies is an array containing different ways of chunking up an amount of cryptocurrency
// for example one strategy is [0.5, 0.5] which represents cutting up the amount into halves
// ideally the strategies would not be predetermined but determined at runtime, but this comes at the cost of
// additional complexity, performance, and rounding errors. so for now only a set number of strategies are used by the tumbler
type strategies map[int][]float64

func getStrategies() strategies {
	s := strategies{
		0: {1},
		1: {0.5, 0.5},
		2: {0.2, 0.2, 0.2, 0.2, 0.2},
		3: {0.4, 0.2, 0.4},
		4: {0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		5: {0.8, 0.2},
	}

	return s
}

func pickRandom(number int) int {
	return rand.Int() % number
}
