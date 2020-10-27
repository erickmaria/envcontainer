package main

import (
	"github.com/ErickMaria/envcontainer/internal/options"
)

var flag options.Flag

func init() {

	flag = options.Flag{}
	flag.Values = make(map[string]string)
	flag.Register("project", "envcontainer", "project name")
	flag.Register("listener", "0", "docker comtainer port listener")
	flag.Register("envfile", "env/.variables", "docker environemt file")
}

func main() {

	command := options.Command{}
	command.Init(flag.Values)

}
