package command

import (
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
)

const (
	EnvCliNoColor = `NO_COLOR`
)

type NamedCommand interface {
	Name() string
}

type CommandFunc func(meta Meta) map[string]cli.CommandFactory

func Commands(metaPtr *Meta, commands CommandFunc) map[string]cli.CommandFactory {
	if metaPtr == nil {
		metaPtr = new(Meta)
	}

	meta := *metaPtr
	if meta.Ui == nil {
		meta.Ui = &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      colorable.NewColorableStdout(),
			ErrorWriter: colorable.NewColorableStderr(),
		}
	}

	all := map[string]cli.CommandFactory{}

	for k, v := range commands(meta) {
		all[k] = v
	}

	return nil
}

type Command interface {
	Name() string
	FlagSet() *flag.FlagSet
	Arguments() []Argument
	Synopsis() string
	Examples() map[string]string
}

func CommandHelp(c Command) string {
	appName := os.Getenv("CLI_APP_NAME")
	helpText := `
Usage: ` + appName + ` ` + c.Name() + ` ` + FlagString(c.FlagSet()) + ` ` + ArgumentAsString(c.Arguments()) + `

	` + c.Synopsis()

	options := c.FlagSet().FlagUsages()
	if options != "" {
		helpText += `
Options:

		` + options
	}

	arguments := ArgumentsString(c.Arguments())
	if arguments != "" {
		helpText += `
Arguments:

		` + arguments
	}

	examples := ExampleString(c.Examples())
	if examples != "" {
		helpText += `
Examples:

` + examples
	}

	return strings.TrimSpace(helpText) + "\n"
}
