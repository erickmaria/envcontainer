package syscmd

import "os"

func DeletePath(path string) error {
	return os.RemoveAll(path)
}
