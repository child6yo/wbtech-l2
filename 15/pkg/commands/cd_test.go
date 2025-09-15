package commands

import (
	"os"
	"testing"
)

func TestChangeDirectorySuccess(t *testing.T) {
	tmpDir := t.TempDir()
	err := ChangeDirectory(tmpDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	currDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	if currDir != tmpDir {
		t.Errorf("expected directory %s, got %s", tmpDir, currDir)
	}
}

func TestChangeDirectoryNonExistent(t *testing.T) {
	err := ChangeDirectory("non-existent-directory-12345")
	if err == nil {
		t.Fatal("expected error for non-existent directory, got nil")
	}
}

func TestChangeDirectoryEmptyPath(t *testing.T) {
	err := ChangeDirectory("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}
