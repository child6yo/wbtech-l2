package app

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"l2.15/pkg/commands"
)

var ErrNoPath = errors.New("path required")
var ErrIntArg = errors.New("integer argument needed")

func StartShell() {
	handleStopSignal()

	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout
	errWriter := os.Stderr

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			} else {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}

		if err = execInput(input, writer); err != nil {
			fmt.Fprintln(errWriter, err)
		}
	}
}

func execInput(input string, writer io.Writer) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	args := strings.Split(input, " ")

	var commands [][]string
	current := []string{}

	for _, arg := range args {
		if arg == "|" {
			if len(current) == 0 {
				return errors.New("syntax error: missing command before |")
			}
			commands = append(commands, current)
			current = []string{}
		} else {
			current = append(current, arg)
		}
	}
	if len(current) == 0 {
		return errors.New("syntax error: trailing |")
	}
	commands = append(commands, current)

	if len(commands) == 1 {
		return runCommand(commands[0], nil, writer)
	}

	var prevOutput string

	for i, cmdArgs := range commands {
		var buf bytes.Buffer

		var stdinReader io.Reader
		if i > 0 {
			stdinReader = strings.NewReader(prevOutput)
		}

		var stdoutWriter io.Writer
		if i == len(commands)-1 {
			stdoutWriter = writer
		} else {
			stdoutWriter = &buf
		}

		err := runCommand(cmdArgs, stdinReader, stdoutWriter)
		if err != nil {
			return err
		}

		if i < len(commands)-1 {
			prevOutput = buf.String()
		}
	}

	return nil
}

func runCommand(args []string, reader io.Reader, writer io.Writer) error {
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return ErrNoPath
		}
		return commands.ChangeDirectory(args[1])
	case "pwd":
		return commands.PrintWorkingDirectory(writer)
	case "echo":
		return commands.PrintArgs(writer, args)
	case "kill":
		pid, err := parseIntArg(args[1])
		if err != nil {
			return err
		}
		return commands.KillProcess(pid)
	case "ps":
		// TODO: for linux
	default:
		cmd := exec.Command(args[0], args[1:]...)

		cmd.Stdin = reader
		cmd.Stderr = os.Stderr
		cmd.Stdout = writer

		return cmd.Run()
	}

	return nil
}

func parseIntArg(arg string) (int, error) {
	pid, err := strconv.Atoi(arg)
	if err != nil {
		return 0, ErrIntArg
	}

	return pid, nil
}

func handleStopSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for range sigChan {
			fmt.Print("> ")
		}
	}()
}
