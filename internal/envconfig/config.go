package envconfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/compose"
	herrors "github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

var (
	envconfigfile = "envcontainer.cfg"
)

type Config struct{}

func (config *Config) Save() {

	if !syscmd.ExistsPath(compose.HOME) {
		herrors.Throw("", errors.New(fmt.Sprintf("cannot found %s folder in this directory", compose.HOME)))
	}

	files := syscmd.ListFiles(compose.HOME)

	var envcontainerFiles string
	var project_name string
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		herrors.Throw("", err)

		if strings.Contains(file, "/.env") {
			envs := strings.Split(string(b), "\n")
			for _, env := range envs {
				if strings.Split(env, "=")[0] == "COMPOSE_PROJECT_NAME" {
					project_name = strings.Split(env, "=")[1]
					break
				}
			}

		}

		envcontainerFiles = envcontainerFiles + fmt.Sprintf("#%s\n%s\n", file, strings.Trim(string(b), "\n"))
	}

	if !syscmd.ExistsPath(Home()) {
		herrors.Throw("", errors.New(fmt.Sprintf("cannot found %s folder in this directory", Home())))
	}

	var saveOrUpdate string = "saved"
	if syscmd.ExistsPath(Home() + "/" + project_name + "/" + envconfigfile) {
		saveOrUpdate = "updated"
	}

	herrors.Throw("cannot create folder", syscmd.CreatePath(Home()+"/"+project_name))
	herrors.Throw("cannot create file", syscmd.CreateFile(Home()+"/"+project_name+"/"+envconfigfile, []byte(envcontainerFiles)))

	fmt.Println("envcontainer configuration was", saveOrUpdate)

}

func (config *Config) List() {
	fmt.Println("ID\tNAME")
	for key, value := range syscmd.ListDir(Home(), false) {
		fmt.Printf("%d\t%s\n", key+1, value)
	}
}

func (config *Config) Get(command *cli.Command) {

	scanner := GetConfig(*command.Flags.Values["name"].ValueString)

	syscmd.DeletePath(compose.HOME)

	err := syscmd.CreatePath(compose.PATH_COMPLETE)
	herrors.Throw("", err)

	var filepath string
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#.envcontainer/") {
			filepath = scanner.Text()[1:len(scanner.Text())]
			continue
		}
		syscmd.AppendFile(filepath, []byte(scanner.Text()+"\n"))
	}

	fmt.Println("envcontainer get configuration done!")

}

func GetConfig(envconfigName string) *bufio.Scanner {

	if envconfigName == "" {
		herrors.Throw("use --name to get configuration", errors.New(""))
	}

	if !syscmd.ExistsPath(Home() + "/" + envconfigName) {
		herrors.Throw(fmt.Sprintf("not found %s folder", envconfigName), errors.New(""))
	}

	envcontainercfg := Home() + "/" + envconfigName + "/" + envconfigfile
	if !syscmd.ExistsPath(envcontainercfg) {
		herrors.Throw(fmt.Sprintf("not found %s file ", envconfigfile), errors.New(""))
	}

	b, err := ioutil.ReadFile(envcontainercfg)
	herrors.Throw("", err)

	return bufio.NewScanner(bytes.NewReader(b))
}

func GetFileDataByMark(scanner *bufio.Scanner, filename string) []byte {

	var data, old, new string
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#.envcontainer/") {
			old = scanner.Text()[1:len(scanner.Text())]
			if strings.HasSuffix(scanner.Text(), filename) {
				new = old
				continue
			}
			if old == new {
				break
			}
		}
		if old == new {
			data = data + scanner.Text() + "\n"
		}
	}

	return []byte(data)
}
