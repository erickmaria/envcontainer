package main

import (
	"github.com/ErickMaria/envcontainer/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.Execute())
}
