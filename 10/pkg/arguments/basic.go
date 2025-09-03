package arguments

import "slices"

// BasicSort выполняет базовую лексикографическую сортировку строк.
func BasicSort(in *[]string) {
	slices.SortFunc(*in, basicSortFunc)
}

func basicSortFunc(a, b string) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}

	return 0
}
