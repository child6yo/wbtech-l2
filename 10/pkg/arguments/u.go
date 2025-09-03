package arguments

import (
	"bufio"
	"fmt"
)

// PrintUniqueOnly пишет в w только уникальные строки.
func PrintUniqueOnly(w *bufio.Writer, output []string) {
	set := make(map[string]struct{})
	for _, el := range output {
		if _, ok := set[el]; !ok {
			fmt.Fprintln(w, el)
		}
		set[el] = struct{}{}
	}
}
