package app

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

// Мок для runCommand
func mockRunCommand(args []string, reader io.Reader, writer io.Writer) error {
	cmd := args[0]
	switch cmd {
	case "echo":
		msg := strings.Join(args[1:], " ")
		writer.Write([]byte(msg))
		return nil
	case "cat":
		if reader != nil {
			io.Copy(writer, reader)
		}
		return nil
	case "fail":
		return errors.New("command failed")
	default:
		return errors.New("unknown command: " + cmd)
	}
}

func TestExecInput(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		if err := execInput("", writer); err != nil {
			t.Errorf("expected nil error for empty input, got %v", err)
		}
		if writer.String() != "" {
			t.Errorf("expected no output, got %q", writer.String())
		}
	})

	t.Run("simple echo", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		if err := execInput("echo hello", writer); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got, want := writer.String(), "hello"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("pipe: echo | cat", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		if err := execInput("echo hello | cat", writer); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got, want := writer.String(), "hello"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("multiple pipes: echo a | echo b | cat", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		if err := execInput("echo a | echo b | cat", writer); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got, want := writer.String(), "b"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("command fails in pipeline", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		err := execInput("echo hello | fail | cat", writer)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed") {
			t.Errorf("expected failure error, got %v", err)
		}
	})

	t.Run("syntax error: missing command before |", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		err := execInput(" | echo hello", writer)
		if err == nil {
			t.Fatal("expected syntax error, got nil")
		}
		if !strings.Contains(err.Error(), "missing command before |") {
			t.Errorf("wrong error message: %v", err)
		}
	})

	t.Run("syntax error: trailing |", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		err := execInput("echo hello |", writer)
		if err == nil {
			t.Fatal("expected syntax error, got nil")
		}
		if !strings.Contains(err.Error(), "trailing |") {
			t.Errorf("wrong error message: %v", err)
		}
	})

	t.Run("leading space and tabs", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		if err := execInput("   echo hello   ", writer); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got, want := writer.String(), "hello"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("single command with multiple args", func(t *testing.T) {
		runCommandFunc = mockRunCommand
		defer func() { runCommandFunc = runCommand }()

		writer := &bytes.Buffer{}
		if err := execInput("echo hello world", writer); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got, want := writer.String(), "hello world"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
