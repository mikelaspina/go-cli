// Copyright 2014 Mike LaSpina. All rights reserved.
// See the LICENSE file for copying permission.

// Package cli extends the command-line parsing provided by the flag package
// with support for subcommands.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

// A RunFunc is a function that invokes a command.
type RunFunc func(args []string) error

// A Command represents an action that can be invoked or a help topic.
type Command struct {
	Run      RunFunc      // non-nil for an invokable action, otherwise a topic
	Usage    string       // usage message
	Short    string       // short (single-line) help text
	Synopsis string       // multi-line help text
	Flags    flag.FlagSet // command line flags
}

// A CommandSet represents a set of named commands.
type CommandSet struct {
	cmds map[string]*Command
	Name string // program name as it should appear in usage; use name() accessor
	Desc string // program description
}

// NewCommandSet creates a new, empty command set.
func NewCommandSet(name, desc string) *CommandSet {
	return &CommandSet{
		cmds: make(map[string]*Command),
		Name: name,
		Desc: desc,
	}
}

// Register adds a named command. Register panics if cmd is nil.
func (cs *CommandSet) Register(name string, cmd *Command) {
	if cmd == nil {
		panic("cli: nil command registered")
	}
	if _, ok := cs.cmds[name]; ok {
		fmt.Fprintf(os.Stderr, "warning: command %q already exits", name)
	}
	if cmd.Flags.Usage == nil {
		cmd.Flags.Usage = func() { cs.PrintUsage(name) }
	}
	cs.cmds[name] = cmd
}

// Run invokes a named command.
func (cs *CommandSet) Run(name string, args []string) error {
	cmd, ok := cs.cmds[name]
	if !ok {
		switch {
		case name != "help":
			cs.unknownCommand(os.Stderr, name)
		case len(args) == 1:
			cs.PrintUsage(args[0])
		default:
			cs.PrintUsage("")
		}

		os.Exit(2)
	}

	if err := cmd.Flags.Parse(args); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(2)
	}

	return cmd.Run(cmd.Flags.Args())
}

// name returns the program name as it should appear in a usage message.
// e.g. name [<options>] <file>...
func (cs *CommandSet) name() string {
	if cs.Name == "" {
		return cs.Name
	}
	return path.Base(os.Args[0])
}

// actions returns a lexicographically sorted list of runnable commands
// in the CommandSet.
func (cs *CommandSet) actions() []string {
	actionNames := make([]string, 0, len(cs.cmds))
	for name := range cs.cmds {
		actionNames = append(actionNames, name)
	}
	sort.Strings(actionNames)
	return actionNames
}

// partialMatch returns an array containing the commands that start
// with prefix in lexicographical order.
func (cs *CommandSet) partialMatch(prefix string) []string {
	var names []string
	for name := range cs.cmds {
		if strings.HasPrefix(name, prefix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func (cs *CommandSet) unknownCommand(w io.Writer, name string) {
	fmt.Fprintf(w, "unknown command: %s\n", name)
	if possibles := cs.partialMatch(name); len(possibles) > 0 {
		fmt.Fprintln(w, "\nDid you mean one of these?")
		for _, s := range possibles {
			fmt.Fprintf(w, "\t%s\n", s)
		}
	}
}

// Default is the default command set.
var Default = NewCommandSet("", "")

// Register adds a named command and panics if cmd is nil.
func Register(name string, cmd *Command) {
	Default.Register(name, cmd)
}

// Run parses the command-line flags from os.Args()[2:], and invokes the
// subcommand named by os.Args()[1].
func Run() error {
	flag.Usage = func() { PrintUsage("") }
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		PrintUsage("")
		os.Exit(2)
	}

	return Default.Run(args[0], args[1:])
}

// PrintUsage prints the usage text for a command to standard error. If
// no command is given, the list of available commands is printed instead.
func PrintUsage(name string) {
	Default.PrintUsage(name)
}
