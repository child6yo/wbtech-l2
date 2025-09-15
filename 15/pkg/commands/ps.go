package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type process struct {
	PID  int
	Name string
}

func getProcesses() ([]process, error) {
	dir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var processes []process
	for _, name := range names {
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		commPath := filepath.Join("/proc", name, "comm")
		data, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}

		cmd := strings.TrimSpace(string(data))
		processes = append(processes, process{PID: pid, Name: cmd})
	}

	return processes, nil
}

func PrintProcesses(writer io.Writer) error {
	return PrintProcessesWithGetter(writer, getProcesses)
}

func PrintProcessesWithGetter(writer io.Writer, getter func() ([]process, error)) error {
	if writer == nil {
		return errors.New("nil input")
	}
	processes, err := getter()
	if err != nil {
		return fmt.Errorf("read processes: %v", err)
	}

	fmt.Fprintf(writer, "%s\t%s\n", "PID", "COMMAND")
	for _, p := range processes {
		fmt.Fprintf(writer, "%d\t%s\n", p.PID, p.Name)
	}
	return nil
}
