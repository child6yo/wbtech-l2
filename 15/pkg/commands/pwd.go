package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// PrintWorkingDirectory печатает текущую рабочую директорию в output.
func PrintWorkingDirectory(writer io.Writer) error {
	if writer == nil {
		return errors.New("nil input")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, dir)

	return err
}
