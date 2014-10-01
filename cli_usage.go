// Copyright 2014 Mike LaSpina. All rights reserved.
// See the LICENSE file for copying permission.

package cli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

func (self *CommandSet) printUsage() {
	eprintf("usage: %s <command> [arguments]\n\n", self.name())

	if names := self.actions(); len(names) > 0 {
		eprintln("Available commands:")
		nameWidth := maxLen(names)
		for _, name := range names {
			eprintf("    %-*s   %s\n", nameWidth, name, self.cmds[name].Short)
		}
		eprintf("\nUse '%s help <command>' for more information on a specific command.\n", self.name())
	}

	eprintln()
}

// cmd.Usage(programName)
func (self *CommandSet) printUsageCmd(cmd *Command) {
	eprintf("usage: %s %s\n\n", self.name(), cmd.Usage)
	eprintf("Arguments:\n")
	columnize(&cmd.Flags)
	if cmd.Synopsis != "" {
		eprintf("\n%s\n", cmd.Synopsis)
	}
}

// columnize aligns a set of flags into two columns. One for the flag
// plus its default value, and one for the usage text. The columns are
// printed to standard output with a left margin of 3 spaces.
func columnize(flags *flag.FlagSet) {
	rows := make([][]string, 0)
	flags.VisitAll(func(f *flag.Flag) {
		rows = append(rows, []string{formatFlag(f), f.Usage})
	})

	flagWidth := 0
	for _, row := range rows {
		if len(row[0]) > flagWidth {
			flagWidth = len(row[0])
		}
	}

	for _, row := range rows {
		eprintf("   %-*s   %s\n", flagWidth, row[0], row[1])
	}
}

func formatFlag(f *flag.Flag) string {
	leading := "-"
	if len(f.Name) > 1 {
		leading = "--"
	}

	format := "%s%s=%s"
	if shouldQuoteValue(f) {
		format = "%s%s=%q"
	}

	return fmt.Sprintf(format, leading, f.Name, f.DefValue)
}

// shouldQuoteValue determines whether a Flag's default value should
// be quoted when printed.
func shouldQuoteValue(f *flag.Flag) bool {
	typ := reflect.TypeOf(f.Value)
	return typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.String
}

func eprintf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

func eprintln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}
