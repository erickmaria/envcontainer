package syscmd

import (
	"os"
)

func ExistsPath(path string) (bool, error) {

	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		return true, nil
	}
	return false, err
}
