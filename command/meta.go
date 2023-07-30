package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/posener/complete"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	shortId = 8
	fullId  = 36
)

type FlagSetFlags uint

const (
	FlagSetNoe     FlagSetFlags = 0
	FlagSetClient  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetClient
)

type Meta struct {
	Ui      cli.Ui
	noColor bool
}

func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	if fs&FlagSetClient != 0 {
		f.BoolVar(&m.noColor, "no-color", false, "disables colored command output. Alternatively, NO_COLOR may be set.")
	}

	f.SetOutput(&uiErrorWriter{ui: m.Ui})

	return f
}

func (m *Meta) AutocompleteFlags(fs FlagSetFlags) complete.Flags {
	if fs&FlagSetClient == 0 {
		return nil
	}

	return complete.Flags{
		"-no-color": complete.PredictNothing,
	}
}

func (m *Meta) Colorize() *colorstring.Colorize {
	return &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: m.noColor || !terminal.IsTerminal(int(os.Stdout.Fd())),
		Reset:   true,
	}
}

type GlobalFlagCommand interface {
	GlobalFlags(*flag.FlagSet)
}

type funcVar func(s string) error

func (f funcVar) Set(s string) error { return f(s) }
func (f funcVar) String() string     { return "" }
func (f funcVar) IsBool() bool       { return false }

func ExampleString(examples map[string]string) string {
	exampleString := []string{}

	for name, example := range examples {
		exampleString = append(exampleString, " "+name+"\n    $ "+example)
	}

	return strings.Join(exampleString, "\n\n")
}

func FlagString(flags *flag.FlagSet) string {
	flagString := []string{}

	flags.VisitAll(func(f *flag.Flag) {
		if f.DefValue == "true" || f.DefValue == "false" {
			flagString = append(flagString, fmt.Sprintf("--%s", f.Name))
			return
		}

		flagString = append(flagString, fmt.Sprintf("--%s <%[1]s-value>", f.Name))
	})
	return strings.Join(flagString, " ")
}
