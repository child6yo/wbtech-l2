package utils

import (
	"bufio"
	"fmt"
)

// PrintResults стандартный вывод результатов во writer.
// Принимает срез строк на вывод (lines), индексы строк, которые необходимо вывести (matches)
// и ряд параметров.
// Выводит только уникальные значения.
func PrintResults(writer *bufio.Writer, lines []string, matches []int,
	beforeContext, afterContext int, showLineNumbers bool) {

	printed := make(map[int]struct{})

	for _, idx := range matches {
		start := max(0, idx-beforeContext)
		end := min(len(lines), idx+afterContext+1)

		for i := start; i < end; i++ {
			if _, ok := printed[i]; ok {
				continue
			}

			output := lines[i]
			if showLineNumbers {
				output = fmt.Sprintf("%d:%s", i+1, output)
			}
			fmt.Fprintln(writer, output)
			printed[i] = struct{}{}
		}
	}
}

// PrintOnlyNumLines выводит только количество совпавших строк во writer.
func PrintOnlyNumLines(writer *bufio.Writer, matches []int) {
	num := len(matches)
	fmt.Fprintln(writer, num)
}
