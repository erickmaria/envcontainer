package options

import (
	"flag"
)

type Values struct {
	value       *string
	Defaulvalue string
	Description string
}

type Flag struct {
	Values map[string]Values
}

func (f *Flag) Register() {

	for k, v := range f.Values {
		v.value = flag.String(k, v.Defaulvalue, v.Description)
		f.Values[k] = v
	}
	flag.Parse()
}
