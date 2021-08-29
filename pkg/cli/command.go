package cli

import (
	"fmt"
	"os"
	"strings"
)

type CommandConfig map[string]Command

type Command struct {
	Flags       Flag
	Quetion     Quetion
	RunAfterAll func()
	Exec        func()
	Desc        string
}

func NewCommand(cc CommandConfig) (*Command, CommandConfig) {

	if isEmpty() {
		Help(cc)
		os.Exit(0)
	}

	if !contains(cc) {
		fmt.Printf(executableName()+": '%s' is not a valid command\n%s\n", os.Args[1], cc["HELP"].Desc)
		os.Exit(0)
	}

	ccAux := cc[os.Args[1]]
	return &ccAux, cc
}

func (c Command) Listener() {
	c.Flags.Register()
	c.RunAfterAll()
	c.Quetion.Make()
	c.Exec()
}

func isEmpty() bool {
	if len(os.Args) < 2 {
		return true
	}
	return false
}

func contains(cc CommandConfig) bool {

	for k := range cc {
		if strings.ToLower(os.Args[1]) == k {
			return true
		}
	}

	return false
}

func Help(descs map[string]Command) {

	fmt.Println("\nUsage: " + executableName() + " COMMAND --FLAGS")

	fmt.Println("\nCommands")

	for commandKey, comandValue := range descs {
		fmt.Printf("%s:     \t%v\n", commandKey, descs[commandKey].Desc)
		for flagKey, flagValue := range comandValue.Flags.Values {
			fmt.Printf("    --%s:     \t%v\n", flagKey, flagValue.Description)
		}

	}

	fmt.Println()
}

func executableName() string {
	executable, err := os.Executable()
	if err != nil {

	}
	return executable
}
