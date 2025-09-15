package commands

import "os"

// KillProcess завершает процесс по переданному PID. 
func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Kill()
}