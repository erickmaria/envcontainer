package main

import (
	"fmt"

	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Version() cli.Command {
	return cli.Command{
		Exec: func() {
			fmt.Println("Version: 2.7.0")
		},
		Desc: "show envcontainer version",
	}
}
