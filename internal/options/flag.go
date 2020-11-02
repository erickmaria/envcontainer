package options

import (
	"flag"
	"os"
)

type Values struct {
	value       *string
	Defaulvalue string
	Description string
}

type Flag struct {
	Command string
	Values  map[string]Values
}

func (f *Flag) Register() {

	cmdFlang := flag.NewFlagSet(f.Command, flag.ExitOnError)

	for k, v := range f.Values {
		v.value = cmdFlang.String(k, v.Defaulvalue, v.Description)
		f.Values[k] = v
	}
	cmdFlang.Parse(os.Args[2:])
}
