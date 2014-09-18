package main

import (
	"fmt"
)

var Version = "dev"

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Display current version",
	Long: `
Display current version

Examples:

	forego version
`,
}

func init() {
}

func runVersion(cmd *Command, args []string) {
	fmt.Println(Version)
}
