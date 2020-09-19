package scan

import (
	"strings"
	"bufio"
	"fmt"
	"log"
	"os"
)

type DockerCompose struct {

}

func (s *DockerCompose) port() (string) {
	return `ports:
      - "$PORT_LISTENER"`
}

func (s *DockerCompose) File(flags map[string]string){
	readFile, err := os.Open("../../configs/docker/compose/docker-compose.yaml")
	
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
 
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string
 
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}
 
	readFile.Close()

	for k, v := range flags {
		fmt.Println(k, v)
	}
	
	for _, eachline := range fileTextLines {
		var myText = strings.Replace(eachline, "${PROJECT}", "Willkommen", -1)
		fmt.Println(myText)
	}
}