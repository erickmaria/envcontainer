package envasync

import (
	"os"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/compose"
	"github.com/ErickMaria/envcontainer/internal/envconfig"
	"github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"github.com/ErickMaria/envcontainer/pkg/cli"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type UpAsyncData struct {
	Compose    []byte
	Dockerfile []byte
	Env        []byte
}

type UpAsync struct {
}

func (UpAsync) Start(command *cli.Command) {
	async := UpAsyncData{}
	var err error

	dockerfileScanner := envconfig.GetConfig(*command.Flags.Values["name"].ValueString)
	async.Dockerfile = envconfig.GetFileDataByMark(dockerfileScanner, "Dockerfile")

	envScanner := envconfig.GetConfig(*command.Flags.Values["name"].ValueString)
	async.Env = envconfig.GetFileDataByMark(envScanner, ".env")

	compose := compose.Compose{}
	dockerComposeScanner := envconfig.GetConfig(*command.Flags.Values["name"].ValueString)
	err = yaml.Unmarshal([]byte(envconfig.GetFileDataByMark(dockerComposeScanner, "docker-compose.yaml")), &compose)
	errors.Throw("", err)

	pwd, err := os.Getwd()
	errors.Throw("", err)

	id := uuid.New()
	envTmp := "/tmp/" + id.String()

	compose.Services.Environment.Volumes[0].Source = pwd
	compose.Services.Environment.EnvFile[0] = envTmp + "/compose/.env"
	compose.Services.Environment.ContainerName = id.String()
	compose.Services.Environment.WorkingDir = "/home/envcontainer/" + strings.Split(pwd, "/")[len(strings.Split(pwd, "/"))-1]

	async.Compose, err = yaml.Marshal(compose)
	errors.Throw("", err)

	syscmd.CreatePath(envTmp + "/compose")
	syscmd.CreateFile(envTmp+"/compose/docker-compose.yaml", async.Compose)
	syscmd.CreateFile(envTmp+"/compose/.env", async.Env)

	syscmd.CreateFile(envTmp+"/Dockerfile", async.Dockerfile)

	build(envTmp)

	syscmd.Exec(
		"docker-compose",
		"-f",
		envTmp+"/compose/docker-compose.yaml",
		"up",
		"-d",
	)

	syscmd.Exec(
		"docker",
		"exec",
		"-it",
		compose.Services.Environment.ContainerName,
		"bash",
	)

	syscmd.Exec(
		"docker-compose",
		"-f",
		envTmp+"/compose/docker-compose.yaml",
		"down",
	)

	errors.Throw("", syscmd.DeletePath(envTmp))

}
func (UpAsync) Exit(command *cli.Command) {

}

func build(basedir string) {

	debug := os.Getenv("ENVCONTAINER_DEBUG")
	if debug != "" && strings.ToLower(debug) == "true" {
		syscmd.Exec(
			"docker-compose",
			"-f",
			basedir+"/compose/docker-compose.yaml",
			"build",
			"-q",
		)
		return
	}

	syscmd.Exec(
		"docker-compose",
		"-f",
		basedir+"/compose/docker-compose.yaml",
		"build",
	)
}
