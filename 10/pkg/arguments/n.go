package arguments

import (
	"slices"
	"strconv"
)

// NumSort выполняет сортровку строк по числовому значению (строки как числа).
func NumSort(in *[]string) {
	slices.SortFunc(*in, numSortFunc)
}

func numSortFunc(a, b string) int {
	aInt, errA := strconv.Atoi(a)
	bInt, errB := strconv.Atoi(b)

	switch {
	case errA == nil && errB != nil:
		return -1
	case errA != nil && errB == nil:
		return 1
	case errA != nil && errB != nil:
		return basicSortFunc(a, b)
	}

	if aInt < bInt {
		return -1
	} else if aInt > bInt {
		return 1
	}
	return 0
}
