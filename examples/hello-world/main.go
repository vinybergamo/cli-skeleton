package main

import (
	"fmt"
	"os"

	"hello-world/commands"

	"github.com/mitchellh/cli"
	"github.com/vinybergamo/cli-skeleton/command"
)

var AppName = "hello-world"

var Version string

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {
	commandMeta := command.SetupRun(AppName, Version, args)
	c := cli.NewCLI(AppName, Version)
	c.Args = os.Args[1:]
	c.Commands = command.Commands(commandMeta, Commands)
	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

func Commands(meta command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"eat": func() (cli.Command, error) {
			return &commands.EatCommand{Meta: meta}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{Meta: meta}, nil
		},
	}
}
