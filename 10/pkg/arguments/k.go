package arguments

import (
	"slices"
	"strings"
)

// ColSort выполняет базовую сортировку по заданной колонке, при разделении
// учитывает только символ табуляции (\t).
// Допролнительно может сортировать по числовому значению при numeric = true.
func ColSort(in *[]string, col int, numeric bool) {
	slices.SortFunc(*in, func(a, b string) int {
		lineA := strings.Split(a, "\t")
		lineB := strings.Split(b, "\t")

		if col >= len(lineA) || col >= len(lineB) {
			return 0
		}

		if numeric {
			return numSortFunc(lineA[col], lineB[col])
		}

		return basicSortFunc(lineA[col], lineB[col])
	})
}
