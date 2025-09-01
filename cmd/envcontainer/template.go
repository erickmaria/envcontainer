package main

import (

	// "text/template"

	envtemplate "github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Template() cli.Command {
	return cli.Command{
		Exec: func() {

			_, _, err := envtemplate.GetConfig(false)
			if err != nil {
				panic(err)
			}

			// funcMap := template.FuncMap{
			// 	"exec": execBashCommand, // Register the function as "exec"
			// }

			// tpl, err := template.New("base").Funcs(sprig.FuncMap()).Funcs(funcMap).Parse(`
			// user: {{ exec "id" "-g" }}
			// currentPath {{ env "" }}
			// `)

			// if err != nil {
			// 	fmt.Println("Error parsing template:", err)
			// 	return
			// }

			// err = tpl.Execute(os.Stdout, nil)
			// if err != nil {
			// 	fmt.Println("Error executing template:", err)
			// }

		},
		Desc: "Testing go templates",
	}
}
