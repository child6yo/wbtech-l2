package commands

import "os"

func ChangeDirectory(dir string) error {
	return os.Chdir(dir)
}
