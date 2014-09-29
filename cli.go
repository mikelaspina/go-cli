// Copyright 2014 Mike LaSpina. All rights reserved.
// See the LICENSE file for copying permission.

package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

// A Runnable is a function invoked to execute a command.
type Runnable func(args []string) error

// A Command represents an action that can be invoked or a help topic.
type Command struct {
	Run   Runnable     // non-nil for an invokable action, otherwise a topic
	Usage string       // usage message
	Short string       // short (single-line) help text
	Long  string       // multi-line help text
	Flags flag.FlagSet // command line flags
}

// A CommandSet represents a set of named commands.
type CommandSet struct {
	cmds map[string]*Command
	Name string // program name as it should appear in usage; use name() accessor
	Desc string // program description
}

// NewCommandSet creates a new, empty command set.
func NewCommandSet() *CommandSet {
	return &CommandSet{
		cmds: make(map[string]*Command),
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
		cmd.Flags.Usage = func() { cs.Help(name) }
	}
	cs.cmds[name] = cmd
}

// Run invokes a named command.
func (cs *CommandSet) Run(name string, args []string) error {
	cmd, ok := cs.cmds[name]
	if !ok {
		if name == "help" {
			if len(args) == 1 {
				cs.Help(args[0])
			} else {
				cs.Usage()
			}
		} else {
			cs.unknownCommand(os.Stderr, name)
		}
		os.Exit(2)
	}

	if err := cmd.Flags.Parse(args); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "%s: %v\n", err)
		}
		os.Exit(2)
	}

	return cmd.Run(cmd.Flags.Args())
}

// Usage prints the help message for the CommandSet to standard error,
// and panics if an error occurs.
func (cs *CommandSet) Usage() {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "usage: %s <command> [arguments]\n\n", cs.name())

	if names := cs.actions(); len(names) > 0 {
		fmt.Fprintln(&buf, "Available commands:")

		nameWidth := maxLen(names)
		for _, name := range names {
			c := cs.cmds[name]
			fmt.Fprintf(&buf, "    %-*s   %s\n", nameWidth, name, c.Short)
		}

		fmt.Fprintf(&buf, "\nUse '%s help <command>' for more information on a specific command.\n", cs.name())
	}

	fmt.Fprintln(&buf)

	if _, err := buf.WriteTo(os.Stderr); err != nil {
		panic(err)
	}
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
	for name, _ := range cs.cmds {
		actionNames = append(actionNames, name)
	}
	sort.Strings(actionNames)
	return actionNames
}

// partialMatch returns an array containing the commands that start
// with prefix in lexicographical order.
func (cs *CommandSet) partialMatch(prefix string) []string {
	names := make([]string, 0)
	for name, _ := range cs.cmds {
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

var defaultCommandSet = NewCommandSet()

// SetName sets the name of the command set.
func SetName(s string) {
	defaultCommandSet.Name = s
}

// SetDescription sets the CommandSet description text.
func SetDescription(s string) {
	defaultCommandSet.Desc = s
}

// Register adds a named command. Register panics if c is nil.
func Register(name string, c *Command) {
	defaultCommandSet.Register(name, c)
}

// Run invokes a named command.
func Run() error {
	flag.Usage = Usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		Usage()
		os.Exit(2)
	}

	return defaultCommandSet.Run(args[0], args[1:])
}

// Usage prints the help message to standard error, and  panics if an
// error occurs.
func Usage() {
	defaultCommandSet.Usage()
}
