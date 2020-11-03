package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type DockerCompose struct {
	Version  string   `yaml:"version"`
	Services Services `yaml:"services"`
}
type Build struct {
	Dockerfile string `yaml:"dockerfile"`
	Context    string `yaml:"context"`
}
type Volumes struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
type Environment struct {
	ContainerName string    `yaml:"container_name"`
	Build         Build     `yaml:"build"`
	Ports         []string  `yaml:"ports"`
	WorkingDir    string    `yaml:"working_dir"`
	EnvFile       []string  `yaml:"env_file"`
	Volumes       []Volumes `yaml:"volumes"`
	Privileged    bool      `yaml:"privileged"`
	StdinOpen     bool      `yaml:"stdin_open"`
	Tty           bool      `yaml:"tty"`
}
type Services struct {
	Environment Environment `yaml:"environment"`
}

var DockerComposeConfig *DockerCompose

func (dc *DockerCompose) Marshal() []byte {

	data, err := yaml.Marshal(dc)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return data
}

func (dc *DockerCompose) Unmarshal(path string) DockerCompose {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	optDc := DockerCompose{}

	err2 := yaml.Unmarshal(data, &optDc)

	if err2 != nil {
		log.Fatalf("error: %v", err)
	}

	return optDc
}

// func (s *DockerCompose) port() (string) {
// 	return `ports:
//       - "$PORT_LISTENER"`
// }

// func (s *DockerCompose) File(flags map[string]string){
// 	readFile, err := os.Open("../../configs/docker/compose/docker-compose.yaml")

// 	if err != nil {
// 		log.Fatalf("failed to open file: %s", err)
// 	}

// 	fileScanner := bufio.NewScanner(readFile)
// 	fileScanner.Split(bufio.ScanLines)
// 	var fileTextLines []string

// 	for fileScanner.Scan() {
// 		fileTextLines = append(fileTextLines, fileScanner.Text())
// 	}

// 	readFile.Close()

// 	for k, v := range flags {
// 		fmt.Println(k, v)
// 	}

// 	for _, eachline := range fileTextLines {
// 		var myText = strings.Replace(eachline, "${PROJECT}", "Willkommen", -1)
// 		fmt.Println(myText)
// 	}
// }
