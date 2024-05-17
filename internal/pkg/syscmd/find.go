package syscmd

import (
	"errors"
	"os"
)

func FindFileCloser(filename string) (string, error) {

	path := "./"

	for i := 1; i < 100; i++ {

		files, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}

		for _, file := range files {
			if filename == file.Name() {
				return path + filename, nil
			}
		}
		
		path = path + "../"
	}

	return "", errors.New("file " + filename + " not found")
}
