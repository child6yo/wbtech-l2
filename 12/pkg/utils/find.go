package utils

// FindMatchesIndexes ищет совпадения функцией match в срезе lines, возвращает индексы совпавших строк.
// При invert = true возвращает несовпадающие строки.
func FindMatchesIndexes(lines []string, match Matcher, invert bool) []int {
	var matches []int
	for i, line := range lines {
		if (match(line) && !invert) || (invert && !match(line)) {
			matches = append(matches, i)
		}
	}
	return matches
}
