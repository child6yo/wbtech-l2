package commands

import (
	"fmt"
	"io"
	"os"
)

func PrintWorkingDirectory(output io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(output, dir)

	return err
}
