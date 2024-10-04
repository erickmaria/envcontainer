package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Input struct {
	Message string
	Default string
}

type Prompt struct {
	Inputs map[string]Input
}

func (quetion *Prompt) Make() {

	var count = 1
	for {

		if count > len(quetion.Inputs) {
			break
		}
		for k, v := range quetion.Inputs {

			if strings.HasPrefix(k, fmt.Sprint(count)) {
				v.Message = quetion.loadRef(v.Message)
				fmt.Print(v.Message)

				reader := bufio.NewReader(os.Stdin)
				output, _, err := reader.ReadLine()
				if err != nil {
					panic(err)
				}
				if len(output) != 0 {
					v.Default = string(output)
				}

				quetion.Inputs[k] = v
				count++
			}
		}
	}

}

func (quetion *Prompt) loadRef(message string) string {
	var result string = message

	re := regexp.MustCompile(`#{{(.*?)}}#`)
	matches := re.FindAllStringSubmatch(message, -1)
	for _, match := range matches {
		result = strings.Replace(result, "#{{"+match[1]+"}}#", quetion.Inputs[match[1]].Default, -1)
	}
	return result
}
