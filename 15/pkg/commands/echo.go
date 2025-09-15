package commands

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// PrintArgs печатает переданные аргументы в output.
func PrintArgs(writer io.Writer, args []string) error {
	if writer == nil || args == nil {
		return errors.New("nil input")
	}
	res := strings.Join(args[1:], " ")
	_, err := fmt.Fprintln(writer, res)

	return err
}
