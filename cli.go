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
	"time"
)

// A RunFunc is a function that invokes a command.
type RunFunc func(args []string) error

// A Command represents an action that can be invoked or a help topic.
type Command struct {
	Run      RunFunc      // non-nil for an invokable action, otherwise a topic
	Usage    string       // usage message
	Short    string       // short (single-line) help text
	Synopsis string       // multi-line help text
	flags    flag.FlagSet // command line flags
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
	if cmd.flags.Usage == nil {
		cmd.flags.Usage = func() { cs.PrintUsage(name) }
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

	if err := cmd.flags.Parse(args); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(2)
	}

	return cmd.Run(cmd.flags.Args())
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

// Bool defines a bool flag with specified name, default value, and usage
// string. The return value is the address of a bool variable that stores
// the value of the flag.
func (cmd *Command) Bool(name string, value bool, usage string) *bool {
	return cmd.flags.Bool(name, value, usage)
}

// BoolVar defines a bool flag with specified name, default value, and usage
// string. The argument p points to a bool variable in which to store the
// value of the flag.
func (cmd *Command) BoolVar(p *bool, name string, value bool, usage string) {
	cmd.flags.BoolVar(p, name, value, usage)
}

// Duration defines a time.Duration flag with specified name, default value,
// and usage string. The return value is the address of a time.Duration
// variable that stores the value of the flag.
func (cmd *Command) Duration(name string, value time.Duration, usage string) *time.Duration {
	return cmd.flags.Duration(name, value, usage)
}

// DurationVar defines a time.Duration flag with specified name, default
// value, and usage string. The argument p points to a time.Duration variable
// in which to store the value of the flag.
func (cmd *Command) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	cmd.flags.DurationVar(p, name, value, usage)
}

// Float64 defines a float64 flag with specified name, default value, and
// usage string. The return value is the address of a float64 variable that
// stores the value of the flag.
func (cmd *Command) Float64(name string, value float64, usage string) *float64 {
	return cmd.Float64(name, value, usage)
}

// Float64Var defines a float64 flag with specified name, default value, and
// usage string. The argument p points to a float64 variable in which to store
// the value of the flag.
func (cmd *Command) Float64Var(p *float64, name string, value float64, usage string) {
	cmd.Float64Var(p, name, value, usage)
}

// Int defines an int flag with specified name, default value, and usage
// string. The return value is the address of an int variable that stores
// the value of the flag.
func (cmd *Command) Int(name string, value int, usage string) *int {
	return cmd.flags.Int(name, value, usage)
}

// IntVar defines an int flag with specified name, default value, and
// usage string. The argument p points to an int variable in which to
// store the value of the flag.
func (cmd *Command) IntVar(p *int, name string, value int, usage string) {
	cmd.flags.IntVar(p, name, value, usage)
}

// Int64 defines an int64 flag with specified name, default value, and
// usage string. The return value is the address of an int64 variable
// that stores the value of the flag.
func (cmd *Command) Int64(name string, value int64, usage string) *int64 {
	return cmd.flags.Int64(name, value, usage)
}

// Int64Var defines an int64 flag with specified name, default value,
// and usage string. The argument p points to an int64 variable in which
// to store the value of the flag.
func (cmd *Command) Int64Var(p *int64, name string, value int64, usage string) {
	cmd.flags.Int64Var(p, name, value, usage)
}

// String defines a string flag with specified name, default value, and
// usage string. The return value is the address of a string variable
// that stores the value of the flag.
func (cmd *Command) String(name string, value string, usage string) *string {
	return cmd.flags.String(name, value, usage)
}

// StringVar defines a string flag with specified name, default value,
// and usage string. The argument p points to a string variable in which
// to store the value of the flag.
func (cmd *Command) StringVar(p *string, name string, value string, usage string) {
	cmd.flags.StringVar(p, name, value, usage)
}

// Uint defines a uint flag with specified name, default value, and usage
// string. The return value is the address of a uint variable that stores
// the value of the flag.
func (cmd *Command) Uint(name string, value uint, usage string) *uint {
	return cmd.flags.Uint(name, value, usage)
}

// UintVar defines a uint flag with specified name, default value, and
// usage string. The argument p points to a uint variable in which to
// store the value of the flag.
func (cmd *Command) UintVar(p *uint, name string, value uint, usage string) {
	cmd.flags.UintVar(p, name, value, usage)
}

// Uint64 defines a uint64 flag with specified name, default value, and
// usage string. The return value is the address of a uint64 variable
// that stores the value of the flag.
func (cmd *Command) Uint64(name string, value uint64, usage string) *uint64 {
	return cmd.flags.Uint64(name, value, usage)
}

// Uint64Var defines a uint64 flag with specified name, default value,
// and usage string. The argument p points to a uint64 variable in which
// to store the value of the flag.
func (cmd *Command) Uint64Var(p *uint64, name string, value uint64, usage string) {
	cmd.flags.Uint64Var(p, name, value, usage)
}

// Var defines a flag with the specified name and usage string.
func (cmd *Command) Var(value flag.Value, name string, usage string) {
	cmd.flags.Var(value, name, usage)
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
