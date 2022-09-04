package file

import "os"

func Link(origin, dest string) error {
	_ = os.RemoveAll(dest)
	return os.Symlink(origin, dest)
}
