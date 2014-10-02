// Copyright 2014 Mike LaSpina. All rights reserved.
// See the LICENSE file for copying permission.

// This example demonstrates use of the cli package.
package cli_test

import "fmt"

// PrintUsage can be used to print a listing of available commands.
func ExampleCommandSet_PrintUsage() {
	version := &cli.Command{
		Usage: "version",
		Short: "print the version and exit",
		Run: func([]string) error {
			fmt.Println("1.0.0")
			return nil
		},
	}

	ui := cli.New("example", "")
	ui.Register("version", version)
	ui.PrintUsage("")
}

// PrintUsage can also be used to print the usage text for a command.
func ExampleCommandSet_PrintUsage_command() {
	export := &cli.Command{
		Usage: "export [-v] [-o <outfile>]",
		Short: "export some data",
		Run: func([]string) error {
			fmt.Println("done!")
			return nil
		},
	}
	export.Bool("-v", false, "cause export to be verbose")
	export.String("-o", "", "output to a file")

	ui := cli.New("my_program", "")
	ui.Register("export", export)
	ui.PrintUsage("export")
}
