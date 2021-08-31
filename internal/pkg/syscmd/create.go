package syscmd

import (
	"io/ioutil"
	"os"

	"github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
)

func CreateFile(name string, data []byte) error {
	return ioutil.WriteFile(name, data, 0644)
}

func AppendFile(name string, data []byte) {
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	errors.Throw("", err)

	if _, err = file.Write(data); err != nil {
		errors.Throw("", err)
	}

}

func CreatePath(path string) error {
	return os.MkdirAll(path, 0755)
}
