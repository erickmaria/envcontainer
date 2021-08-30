package envconfig

import "os"

// var (
// 	HOME = "/.envconfig"
// )

func Home() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homedir + "/.envconfig"
}
