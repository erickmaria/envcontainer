package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type CommandConfig map[string]Command

type Command struct {
	Flags        Flag
	Quetion      Quetion
	RunBeforeAll func()
	Exec         func()
	Desc         string
}

func NewCommand(cc CommandConfig) (*Command, CommandConfig) {

	if isEmpty() {
		Help(cc)
		os.Exit(0)
	}

	if !contains(cc) {
		fmt.Printf(ExecutableName()+": '%s' is not a valid command\n%s\n", os.Args[1], cc["HELP"].Desc)
		os.Exit(0)
	}

	ccAux := cc[os.Args[1]]
	return &ccAux, cc
}

func (c Command) Listener() {
	c.Flags.Register()
	if c.RunBeforeAll != nil {
		c.RunBeforeAll()
	}
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

	fmt.Println("\nUsage: " + ExecutableName() + " COMMAND --FLAGS")

	fmt.Println("\nCommands")

	sortDescs := make([]string, 0, len(descs))
	for k := range descs {
		sortDescs = append(sortDescs, k)
	}

	sort.Strings(sortDescs)
	for _, v := range sortDescs {
		fmt.Printf("%s:     \t\t%v\n", v, descs[v].Desc)
		for flagKey, flagValue := range descs[v].Flags.Values {
			fmt.Printf("    --%s:     \t\t%v\n", flagKey, flagValue.Description)
		}
	}

	// for commandKey, comandValue := range descs {
	// 	fmt.Printf("%s:     \t\t%v\n", commandKey, descs[commandKey].Desc)
	// 	for flagKey, flagValue := range comandValue.Flags.Values {
	// 		fmt.Printf("    --%s:     \t\t%v\n", flagKey, flagValue.Description)
	// 	}

	// }

	fmt.Println()
}

func ExecutableName() string {
	executable, err := os.Executable()
	if err != nil {

	}
	executableNameSplit := strings.Split(executable, "/")
	return executableNameSplit[len(executableNameSplit)-1]
}
