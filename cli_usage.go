// Copyright 2014 Mike LaSpina. All rights reserved.
// See the LICENSE file for copying permission.

package cli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

// PrintUsage prints the usage text for a command to standard error. If
// no command is given, the list of available commands is printed instead.
func (cs *CommandSet) PrintUsage(name string) {
	if cmd, ok := cs.cmds[name]; ok {
		cs.printUsageCmd(cmd)
		return
	}

	eprintf("usage: %s <command> [arguments]\n\n", cs.name())

	if names := cs.actions(); len(names) > 0 {
		eprintln("Available commands:")
		nameWidth := maxLen(names)
		for _, name := range names {
			eprintf("    %-*s   %s\n", nameWidth, name, cs.cmds[name].Short)
		}
		eprintf("\nUse '%s help <command>' for more information on a specific command.\n", cs.name())
	}

	eprintln()
}

func maxLen(ary []string) int {
	max := 0
	for _, s := range ary {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

// cmd.Usage(programName)
func (cs *CommandSet) printUsageCmd(cmd *Command) {
	eprintf("usage: %s %s\n\n", cs.name(), cmd.Usage)
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
	var rows [][2]string
	flags.VisitAll(func(f *flag.Flag) {
		rows = append(rows, [2]string{formatFlag(f), f.Usage})
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
