package syscmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
)

func Exec(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	debug := os.Getenv("ENVCONTAINER_DEBUG")
	if debug != "" && strings.ToLower(debug) == "true" {
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()

	errors.Throw("command failed, check envcontainer configs.", err)
}
