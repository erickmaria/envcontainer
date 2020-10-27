package options

import (
	"flag"
)

type Flag struct {
	Values map[string]string
}

func (f *Flag) Register(flagname, defaulvalue, description string){
	var value string

	flag.StringVar(&value, flagname, defaulvalue, description)
	// flag.StringVar(&value, flagshortname, "", "project name")
	flag.Parse()

	f.Values[flagname] = value
	// f.Values[flagshortname] = value	
}