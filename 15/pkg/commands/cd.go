package commands

import "os"

// ChangeDirectory меняет текущую директорию на переданную.
func ChangeDirectory(dir string) error {
	return os.Chdir(dir)
}
