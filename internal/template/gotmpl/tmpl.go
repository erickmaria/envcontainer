package gotmpl

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"strings"
)

func FuncMap() map[string]any {
	return template.FuncMap{
		"exec": Exec, // Register the function as "exec"
	}
}

func Exec(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing command: %w", err)
	}
	result := out.String()
	return strings.TrimRight(result, "\n"), nil
}
