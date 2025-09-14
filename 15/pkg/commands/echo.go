package commands

import (
	"fmt"
	"io"
	"strings"
)

func PrintArgs(output io.Writer, args []string) error {
	res := strings.Join(args[1:], " ")
	_, err := fmt.Fprintln(output, res)

	return err
}
