package pkg

import (
	"os"
)

func CreateFile(files []string){
	
	for _, file := range files {
		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

}

