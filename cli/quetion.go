package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Query struct {
	Scene string
	Value string
}

type Quetion struct {
	Queries map[string]Query
}

func (quetion *Quetion) Make() {

	var count = 1
	for {

		if count > len(quetion.Queries) {
			break
		}
		for k, v := range quetion.Queries {

			if strings.HasPrefix(k, fmt.Sprint(count)) {
				fmt.Print(v.Scene)
				reader := bufio.NewReader(os.Stdin)
				output, _, err := reader.ReadLine()
				if err != nil {
					panic(err)
				}
				if len(output) != 0 {
					v.Value = string(output)
				}

				quetion.Queries[k] = v
				count++
			}
		}
	}

}
