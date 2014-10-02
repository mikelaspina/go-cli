go-cli
======

A simple library on top of the flag package for building subcommand-based interfaces.

## Getting Started

~~~~
go get github.com/mikelaspina/go-cli
~~~~

```go
package main

import (
	"fmt"
	
	"github.com/mikelaspina/go-cli"
)

type sendCommand struct {
	cli.Command,
	recipient string,
	sender    string,
	message   string,
	verbose   bool
}

func init() {
	send := &sendCommand{
		Usage: "send --from=<email> --to=<email> [-v] message...",
		Short: "send a mail message"
	}
	send.Run = func(args []string) error {
		send.message = strings.Join(args, " ")
		return send.Send()
	}
	send.StringVar(&send.sender, "from", "", "sender email address")
	send.StringVar(&send.message, "to", "", "recipient email address")
	send.BoolVar(&send.verbose, "v", false, "verbose output")
	cli.Register("send", &send.Command)
}

func main() {
    versionCommand := &cli.Command{
		Usage: "version",
		Short: "show version and exit",
		Run:   func([]string) error {
			fmt.Println("1.0.0")
			return nil
		},
	}
	cli.Register("version", versionCommand)

	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (s *sendCommand) Send() error {
	// TODO: validate args and send the message
}
```

## License

This project is released under the [MIT License](http://www.opensource.org/licenses/MIT).
