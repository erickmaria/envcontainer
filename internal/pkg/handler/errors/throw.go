package errors

import (
	"fmt"
	"os"
)

const (
	errorColor = "\033[1;31m%s\033[0m"
)

func Throw(message string, e error) {
	if e != nil {
		fmt.Print(message)

		if e.Error() != "" {
			fmt.Println()
		}
		fmt.Println(e.Error())
		os.Exit(0)
	}
}
