package commands

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestKillProcessRunningChild(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test not supported on Windows")
	}

	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	pid := cmd.Process.Pid
	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	time.Sleep(100 * time.Millisecond)

	err = KillProcess(pid)
	if err != nil {
		t.Fatalf("failed to kill process: %v", err)
	}

	select {
	case waitErr := <-done:
		if waitErr == nil {
			t.Fatal("expected process to exit due to kill, but no error from Wait")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("process did not terminate after kill")
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		t.Fatalf("failed to find process by PID: %v", err)
	}

	err = proc.Release()
	if err != nil {
		t.Fatalf("failed to release process handle: %v", err)
	}
}

func TestKillProcessAlreadyExited(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test not supported on Windows")
	}

	cmd := exec.Command("true")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	err = cmd.Wait()
	if err != nil && err.Error() == "wait: no child processes" {
		t.Skip("child process already gone, cannot test")
	}

	pid := cmd.Process.Pid

	err = KillProcess(pid)
	if err == nil {
		t.Fatal("expected error when killing already exited process, got nil")
	}
}

func TestKillProcessInvalidPID(t *testing.T) {
	err := KillProcess(-1)
	if err == nil {
		t.Fatal("expected error for invalid PID, got nil")
	}
}
