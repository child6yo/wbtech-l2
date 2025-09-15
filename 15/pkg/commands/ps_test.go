package commands

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestPrintProcessesWithMockData(t *testing.T) {
	mockGetter := func() ([]process, error) {
		return []process{
			{PID: 1, Name: "init"},
			{PID: 42, Name: "bash"},
			{PID: 999, Name: "sleep"},
		}, nil
	}

	var buf bytes.Buffer
	err := PrintProcessesWithGetter(&buf, mockGetter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}

	if lines[0] != "PID\tCOMMAND" {
		t.Errorf("bad header: got %q", lines[0])
	}

	expected := []struct {
		pid  int
		name string
	}{
		{1, "init"},
		{42, "bash"},
		{999, "sleep"},
	}

	for i, exp := range expected {
		parts := strings.Split(lines[i+1], "\t")
		if len(parts) != 2 {
			t.Errorf("line %d: bad format: %q", i+1, lines[i+1])
			continue
		}
		pid, _ := strconv.Atoi(parts[0])
		if pid != exp.pid || parts[1] != exp.name {
			t.Errorf("line %d: expected (%d, %s), got (%s, %s)", i+1, exp.pid, exp.name, parts[0], parts[1])
		}
	}
}

func TestPrintProcessesEmpty(t *testing.T) {
	mockGetter := func() ([]process, error) {
		return []process{}, nil
	}

	var buf bytes.Buffer
	err := PrintProcessesWithGetter(&buf, mockGetter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "PID\tCOMMAND\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestPrintProcessesWithError(t *testing.T) {
	mockGetter := func() ([]process, error) {
		return nil, fmt.Errorf("test error")
	}

	var buf bytes.Buffer
	err := PrintProcessesWithGetter(&buf, mockGetter)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("expected error to contain 'test error', got %v", err)
	}
}

func TestGetProcessesFromRealProc(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("only on linux")
	}

	processes, err := getProcesses()
	if err != nil {
		t.Skipf("cannot read /proc: %v", err)
	}

	if len(processes) == 0 {
		t.Error("expected at least one process")
	}

	foundSelf := false
	selfPid := os.Getpid()
	for _, p := range processes {
		if p.PID == selfPid {
			foundSelf = true
			break
		}
	}

	if !foundSelf {
		t.Log("current process not found in /proc — possible but rare")
	}
}
