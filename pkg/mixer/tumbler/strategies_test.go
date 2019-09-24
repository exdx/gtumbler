package tumbler

import (
	"reflect"
	"testing"
)

func TestTumbler_GetStrategies(t *testing.T) {
	expected := &strategies{
		0: {1},
		1: {0.5, 0.5},
		2: {0.2, 0.2, 0.2, 0.2, 0.2},
		3: {0.4, 0.2, 0.4},
		4: {0.8, 0.2},
	}
	got := getStrategies()

	if !reflect.DeepEqual(*expected, *got) {
		t.Errorf("expected strategies %v got %v", *expected, *got)
	}
}

func TestTumbler_PickRandom(t *testing.T) {
	number := 5
	expected := [5]int{0, 1, 2, 3, 4}
	got := pickRandom(number)

	for _, random := range expected {
		if random != got {
			continue
		} else {
			return
		}
	}

	t.Errorf("expected a number from %v, received %d", expected, got)
}
