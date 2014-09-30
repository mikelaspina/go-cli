// Copyright 2014 Mike LaSpina. All rights reserved.
// See the LICENSE file for copying permission.

package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
)

// Help prints the usage text for a command to standard error.
// It panics if an error occurs.
func (cs *CommandSet) Help(name string) {
	buf := bytes.Buffer{}
	if cmd, ok := cs.cmds[name]; ok {
		fmt.Fprintf(&buf, "usage: %s %s\n\n", cs.name(), cmd.Usage)

		buf.WriteString("Arguments:\n")
		columnizeFlags(&buf, &cmd.Flags)

		if cmd.Synopsis != "" {
			buf.WriteString("\n")
			buf.WriteString(cmd.Synopsis)
			buf.WriteString("\n")
		}
	} else {
		cs.unknownCommand(&buf, name)
	}

	if _, err := buf.WriteTo(os.Stderr); err != nil {
		panic(err)
	}
}

func columnizeFlags(w io.Writer, flags *flag.FlagSet) {
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
		fmt.Fprintf(w, "   %-*s   %s\n", flagWidth, row[0], row[1])
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
