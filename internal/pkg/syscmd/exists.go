package syscmd

import (
	"os"
)

func ExistsPath(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
