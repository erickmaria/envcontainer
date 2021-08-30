package syscmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
)

func list(path string, onlyfile bool) []string {

	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if onlyfile && !info.IsDir() {
			files = append(files, path)
		}
		return nil

	})
	errors.Throw("", err)

	return files
}

func ListFiles(path string) []string {
	return list(path, true)
}

func ListDir(path string, subdir bool) []string {
	if !subdir {
		fileInfo, err := ioutil.ReadDir(path)
		errors.Throw("", err)

		var files []string
		for _, file := range fileInfo {
			files = append(files, file.Name())
		}

		return files

	}
	return list(path, false)
}
