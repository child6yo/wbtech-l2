package commands

import (
	"bytes"
	"testing"
)

func TestPrintArgsSingleArg(t *testing.T) {
	var buf bytes.Buffer
	args := []string{"echo", "hello"}
	err := PrintArgs(&buf, args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "hello\n"
	got := buf.String()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPrintArgsMultipleArgs(t *testing.T) {
	var buf bytes.Buffer
	args := []string{"echo", "hello", "world", "go"}
	err := PrintArgs(&buf, args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "hello world go\n"
	got := buf.String()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPrintArgsNoArgs(t *testing.T) {
	var buf bytes.Buffer
	args := []string{"echo"}
	err := PrintArgs(&buf, args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "\n"
	got := buf.String()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPrintArgsWithSpaces(t *testing.T) {
	var buf bytes.Buffer
	args := []string{"echo", "a", "b c", "d"}
	err := PrintArgs(&buf, args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "a b c d\n"
	got := buf.String()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPrintArgsNilWriter(t *testing.T) {
	args := []string{"echo", "test"}
	err := PrintArgs(nil, args)
	if err == nil {
		t.Fatal("expected error when writer is nil, got nil")
	}
}
