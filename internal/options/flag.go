package options

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var getBools = map[string]string{}

type Values struct {
	valueString *string
	valueBool   *bool
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

		if pased, err := strconv.ParseBool(v.Defaulvalue); err == nil {
			v.valueBool = cmdFlang.Bool(k, pased, v.Description)

			f.Values[k] = v
			getBools[k] = v.Defaulvalue
			continue
		}

		v.valueString = cmdFlang.String(k, v.Defaulvalue, v.Description)
		f.Values[k] = v
	}

	cmdFlang.Parse(os.Args[2:])

	UsageMsg := func(notArg string) {
		fmt.Printf("envcontainer: '%s' not is argument\n", notArg)
		cmdFlang.Usage()
		os.Exit(0)
	}

	for k, arg := range os.Args[2:] {

		// check if is key begins with '--' or '-'
		if checkIsArg(arg, false) {

			// check if next value is a arg and if this arg is boolean or not
			if len(os.Args) != k+3 {
				if checkIsArg(os.Args[k+3], false) {
					if !checkIsArg(arg, true) {
						fmt.Printf("envcontainer: '%s' the value is required\n", arg)
						os.Exit(0)
					}
				}
			}

			// check if value is a boolean
			if tryParse(getBools[arg[2:]]) {
				// check is array ends
				if len(os.Args) == k+3 {
					continue
				}
				// check if boolean key have value
				if !checkIsArg(os.Args[k+3], false) {
					// if have a value show usage
					UsageMsg(os.Args[k+3])
				}
				continue
			}

			if len(os.Args) != k+3 {
				if len(os.Args) != k+4 {
					if !checkIsArg(os.Args[k+4], false) {
						UsageMsg(os.Args[k+4])
					}
				}
			}
		}
	}

	cmdFlang.Parse(os.Args[2:])
}

func tryParse(strg string) bool {
	if _, err := strconv.ParseBool(strg); err == nil {
		return true
	}
	return false
}

func checkIsArg(arg string, checkValueToo bool) bool {

	if arg[:2] == "--" {
		if checkValueToo {
			return tryParse(getBools[arg[2:]])
		}
		return true
	} else if arg[:1] == "-" {
		if checkValueToo {
			return tryParse(getBools[arg[1:]])
		}
		return true
	}
	return false
}
