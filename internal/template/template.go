package template

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"github.com/ErickMaria/envcontainer/internal/pkg/types"
	tplTypes "github.com/ErickMaria/envcontainer/internal/template/types"

	"github.com/ErickMaria/envcontainer/internal/template/gotmpl"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

var (
	paths = map[string]string{
		"home":        "/tmp/.envcontainer",
		"dockerfiles": "/tmp/.envcontainer/dockerfiles",
	}
	fileLocation string = ".envcontainer.yaml"
)

type Type string

// Type constants
const (
	TypeBind   Type = "bind"
	TypeVolume Type = "volume"
)

func NewConfigFile(path string) error {

	dest := filepath.Join(path, ".envcontainer.yaml")

	if _, err := os.Stat(dest); err == nil {
		fmt.Printf("%s already exists\n", dest)
		return nil
	}
	// Build a default Envcontainer following the struct in internal/template
	env := tplTypes.Envcontainer{
		Project: tplTypes.Project{
			Name:        filepath.Base(path),
			Version:     "0.0.1",
			Description: "A short description of the project",
		},
		Container: tplTypes.Container{
			Shell:       "bash",
			Ports:       []string{},
			Build:       "FROM alpine:3.18\nRUN apk add --no-cache bash\n",
			NetworkMode: "bridge",
			Networks:    []types.Network{},
		},
		Mounts: []types.Mount{
			{
				Source:   "./",
				Target:   "/workdir",
				Type:     "bind",
				Readonly: false,
			},
		},
		AlwaysUpdate: false,
		AutoStop:     false,
	}

	// Marshal only the project section; optional fields will be added as commented notes
	proj := map[string]any{
		"project": map[string]string{
			"name":        env.Project.Name,
			"version":     env.Project.Version,
			"description": env.Project.Description,
		},
	}

	projBytes, err := yaml.Marshal(proj)
	if err != nil {
		fmt.Printf("failed to marshal project config: %v\n", err)
		return err
	}

	// prepare build block indented
	buildIndented := ""
	if env.Container.Build != "" {
		// indent each line by 4 spaces
		buildIndented = "    " + strings.ReplaceAll(strings.TrimRight(env.Container.Build, "\n"), "\n", "\n    ") + "\n"
	}

	// Assemble content: project (real) + container.build (real) + commented optional fields with examples
	content := string(projBytes) + "\n" +
		"container:\n" +
		"  # shell: bash  # Optional: shell to use inside the container (default: bash)\n" +
		"  # ports:\n" +
		"  #   - \"8080:80\"  # Optional: host:container port mapping\n"

	if buildIndented != "" {
		content += "  build: |\n" + buildIndented + "\n"
	} else {
		content += "  # build: |\n  #   FROM alpine:3.18\n  #   RUN apk add --no-cache bash\n\n"
	}

	content += "  # network_mode: \"host\"  # Optional: Docker network_mode\n" +
		"  # networks:\n" +
		"  #   - name: mynet\n" +
		"  #     external: true\n\n" +
		"# always_update: false  # Optional: pull/update image before starting\n" +
		"# auto_stop: false      # Optional: stop container after run\n\n" +
		"# mounts:  # Optional: bind/volume mounts\n" +
		"# - type: bind\n" +
		"#   source: ./\n" +
		"#   target: /workdir\n" +
		"#   readonly: false\n"

	if err := os.WriteFile(dest, []byte(content), 0644); err != nil {
		fmt.Printf("failed to write %s: %v\n", dest, err)
		return err
	}

	return nil
}

func Initialization() error {

	err := syscmd.CreateDir(toSlice(paths))
	if err != nil {
		return err
	}

	return nil
}

func Unmarshal() (tplTypes.Envcontainer, error) {

	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return tplTypes.Envcontainer{}, err
	}

	var envcontainer tplTypes.Envcontainer
	err = yaml.Unmarshal(data, &envcontainer)
	if err != nil {
		return tplTypes.Envcontainer{}, err
	}

	envcontainer.Project.Name = strings.ReplaceAll(strings.ToLower(envcontainer.Project.Name), " ", "-")
	envcontainer.Container.Build, err = tmpDockerfile(envcontainer)

	if err != nil {
		return tplTypes.Envcontainer{}, err
	}

	err = validate(envcontainer)

	if err != nil {
		return tplTypes.Envcontainer{}, err
	}

	return envcontainer, nil
}

func UnmarshalWithFile(location string) (tplTypes.Envcontainer, error) {

	fileLocation = location

	return Unmarshal()
}

func validate(envcontainer tplTypes.Envcontainer) error {

	if len(envcontainer.Container.Ports) > 0 {

		regex := regexp.MustCompile(`^(\d+)(:?)(\d+)$`)

		for _, v := range envcontainer.Container.Ports {

			if ok := regex.MatchString(v); !ok {
				return errors.New("port " + v + " is invalid")
			}
		}
	}

	return nil
}

func tmpDockerfile(envcontainer tplTypes.Envcontainer) (string, error) {

	_, err := os.ReadFile(envcontainer.Container.Build)
	if err != nil {
		dockerfilePath := GetTmpDockerfileDir(envcontainer)

		err = syscmd.CreateDir([]string{dockerfilePath})
		if err != nil {
			return "", err
		}
		dockerfile := dockerfilePath + "/Dockerfile"
		err = syscmd.CreateFile(dockerfile, []byte(envcontainer.Container.Build))
		if err != nil {
			return "", err
		}

		processDockerfileTemplate(dockerfile)

		return dockerfile, nil

	}
	return envcontainer.Container.Build, nil
}

func processDockerfileTemplate(dockerfile string) {
	tpl, err := template.New(filepath.Base(dockerfile)).
		Funcs(sprig.FuncMap()).
		Funcs(gotmpl.FuncMap()).
		ParseFiles(dockerfile)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Create a buffer to hold the processed template
	var buf strings.Builder
	err = tpl.ExecuteTemplate(&buf, filepath.Base(dockerfile), nil)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// Write the processed template back to the same file
	err = os.WriteFile(dockerfile, []byte(buf.String()), 0644)
	if err != nil {
		fmt.Println("Error writing processed template:", err)
	}
}

func GetTmpDockerfileDir(envcontainer tplTypes.Envcontainer) string {
	return paths["dockerfiles"] + "/" + envcontainer.Project.Name + "/" + envcontainer.Project.Version
}

func List() (map[string]tplTypes.Envcontainer, error) {

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	root := usr.HomeDir
	pattern := ".envcontainer.yaml"

	var matches []string

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == pattern {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var envcontainers = map[string]tplTypes.Envcontainer{}
	for _, match := range matches {
		envcontainer, err := UnmarshalWithFile(match)
		if err != nil {
			fmt.Println(match)
			return nil, err
		}
		envcontainers[match] = envcontainer

	}

	return envcontainers, nil
}

func GetConfig(getCloser bool) (tplTypes.Envcontainer, string, error) {

	configFile, errConfigFile := Unmarshal()
	var defaultMountDir string

	if getCloser {
		file, err := syscmd.FindFileCloser(".envcontainer.yaml")
		if err != nil {
			return tplTypes.Envcontainer{}, "", err
		}

		pwd, _ := os.Getwd()
		for i := 0; i < strings.Count(file, "../"); i++ {
			pwd = strings.Join(strings.Split(pwd, "/")[:len(strings.Split(pwd, "/"))-1], "/")

		}

		if file != "" {
			configFile, err = UnmarshalWithFile(file)
			if err != nil {
				return tplTypes.Envcontainer{}, "", err
			}

		}

		defaultMountDir = pwd + "/.envcontainer/"

	} else if errConfigFile != nil {
		return tplTypes.Envcontainer{}, "", errConfigFile
	}

	if configFile.Container.Shell == "" {
		configFile.Container.Shell = "bash"
	}

	return configFile, defaultMountDir, nil

}

func toSlice(maps map[string]string) []string {

	values := []string{}
	for _, v := range maps {
		values = append(values, v)
	}

	return values
}

// func sliceDeleteEmpty(s []string) []string {
// 	var r []string
// 	for _, str := range s {
// 		if str != "" {
// 			r = append(r, str)
// 		}
// 	}
// 	return r
// }
