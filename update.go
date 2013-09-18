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
	forego update 0.3.0
`,
}

func init() {
}

func runUpdate(cmd *Command, args []string) {
	d := dist.NewDist("https://godist.herokuapp.com", "ddollar/forego")
	if len(args) > 0 {
		err := d.UpdateTo(args[0])
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else {
			fmt.Printf("updated to %s\n", args[0])
		}
	} else {
		version, err := d.Update()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else {
			fmt.Printf("updated to %s\n", version)
		}
	}
}
