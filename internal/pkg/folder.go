package pkg

import (
	"os"
)


func CreateDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				panic(err)
			}
	}
}