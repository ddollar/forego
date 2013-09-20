package main

import (
	"fmt"
	"github.com/ddollar/dist"
)

var cmdUpdate = &Command{
	Run:   runUpdate,
	Usage: "update [version]",
	Short: "Update forego",
	Long: `
Update forego

Examples:

	forego update
	forego update 0.7.0
`,
}

func init() {
}

func runUpdate(cmd *Command, args []string) {
	if Version == "dev" {
		fmt.Println("ERROR: can't update dev version")
		return
	}
	d := dist.NewDist("ddollar/forego")
	var err error
	var to string
	if len(args) > 0 {
		err = d.UpdateTo(Version, args[0])
		to = args[0]
	} else {
		to, err = d.Update(Version)
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	} else {
		fmt.Printf("updated to %s\n", to)
	}
}
