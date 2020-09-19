package pkg

import (
	"os"
)

func WriteFile(file ,text string){

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	f.WriteString(text)
	f.Close()

}

