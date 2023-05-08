package syscmd

import (
	"os"
	// "github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
)

func CreateFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}

func AppendFile(name string, data []byte) error {
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	if _, err = file.Write(data); err != nil {
		return err
	}

	return nil
}

func CreateDir(paths []string) error {

	var err error
	for _, path := range paths {
		if err = os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}
