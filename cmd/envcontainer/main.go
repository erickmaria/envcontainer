package main

import (
	"github.com/ErickMaria/envcontainer/internal/options"
)

var flag options.Flag

func init() {

	flag = options.Flag{
		Values: map[string]options.Values{
			"project": options.Values{
				Defaulvalue: "app",
				Description: "project name",
			},
			"listener": options.Values{
				Defaulvalue: "",
				Description: "docker comtainer port listener",
			},
			"envfile": options.Values{
				Defaulvalue: "env/.variables",
				Description: "docker environemt file",
			},
			"image": options.Values{
				Defaulvalue: "ubuntu",
				Description: "dockerfile image",
			},
		},
	}
	flag.Register()

}

func main() {

	command := options.Command{}
	command.Init(flag)

}
