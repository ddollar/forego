package main

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

func runVersion(cmd *Command, args []string) {
	Println(Version)
}
