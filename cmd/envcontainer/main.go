package main

import (
	"fmt"
	"github.com/ErickMaria/envcontainer/internal/options"
)

var flag options.Flag

func init() {

	flag = options.Flag{}
	flag.Values = make(map[string]string)
	flag.Register("project", "project name")
	flag.Register("listener", "docker comtainer port listener")
	flag.Register("envfile", "docker environemt file")
}

func main() {

	command := options.Command{}
	command.Init(flag.Values)

	fmt.Print()
}
