package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	ttmpt "text/template"

	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Init() cli.Command {

	dir, _ := os.Getwd()
	projectName := strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1]

	return cli.Command{
		Desc: "create a envcontainer.yaml file",
		Flags: cli.Flag{
			Values: map[string]cli.Values{
				"build": {
					Defaulvalue: "false",
					Description: "build envcontainer image",
				},
				"override": {
					Defaulvalue: "false",
					Description: "override envcontainer configs",
				},
			},
		},
		Prompt: cli.Prompt{
			Inputs: map[string]cli.Input{
				"1_project_name": {
					Message: "What's your Project Name [" + projectName + "]: ",
					Default: projectName,
				},
				"2_project_description": {
					Message: "What's your Project Description [Create a development environment for #{{1_project_name}}# ]: ",
					Default: "Create a development environment for " + projectName,
				},
				"3_project_version": {
					Message: "What's your Project Version [0.0.1]: ",
					Default: "0.0.1",
				},
				"4_container_base_image": {
					Message: "What's your Container Base Image [ubuntu:latest]: ",
					Default: "ubuntu:latest",
				},
			},
		},
		RunBeforeAll: func() {

			override := *cmd.Flags.Values["override"].ValueBool

			if !override {
				_, _, err := template.GetConfig(false)
				if err == nil {
					fmt.Println("Envcontainer already initialized in this project! use flag --override to replace current setting.")
					os.Exit(0)
				}
			}
		},
		Exec: func() {

			BuildTpl := `FROM {{.}}
`
			t, err := ttmpt.New("build").Parse(BuildTpl)
			if err != nil {
				panic(err)
			}

			var output bytes.Buffer
			err = t.Execute(&output, cmd.Prompt.Inputs["4_container_base_image"].Default)
			if err != nil {
				panic(err)
			}

			tpl := template.Envcontainer{
				Project: template.Project{
					Name:        cmd.Prompt.Inputs["1_project_name"].Default,
					Description: cmd.Prompt.Inputs["2_project_description"].Default,
					Version:     cmd.Prompt.Inputs["3_project_version"].Default,
				},
				Container: template.Container{
					Build: output.String(),
				},
			}

			err = template.Marshal(tpl)
			if err != nil {
				panic(err)
			}

			fmt.Println("\nEnvcontainer initialized successfully!")
		},
	}
}
