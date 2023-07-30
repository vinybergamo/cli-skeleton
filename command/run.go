package command

import (
	"os"

	colorable "github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func SetupRun(appName string, version string, args []string) *Meta {
	args = SetupEnv(args)

	metaPtr := new(Meta)

	color := true
	if os.Getenv(EnvCliNoColor) != "" {
		color = false
	}

	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))

	if isTerminal && color {
		metaPtr.Ui = &cli.ConcurrentUi{
			Ui: &cli.ColoredUi{
				ErrorColor: cli.UiColorRed,
				WarnColor:  cli.UiColorYellow,
				InfoColor:  cli.UiColorGreen,
				Ui: &cli.BasicUi{
					Reader:      os.Stdin,
					Writer:      colorable.NewColorableStdout(),
					ErrorWriter: colorable.NewColorableStderr(),
				},
			},
		}
	} else {
		metaPtr.Ui = &cli.ConcurrentUi{
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      colorable.NewColorableStdout(),
				ErrorWriter: colorable.NewColorableStderr(),
			},
		}
	}

	os.Setenv("CLI_APP_NAME", appName)
	os.Setenv("CLI_VERSION", version)

	return metaPtr
}

func SetupEnv(args []string) []string {
	noColor := false
	for _, arg := range args {
		if arg == "-no-color" || arg == "--no-color" {
			noColor = true
		}
	}

	if noColor {
		os.Setenv(EnvCliNoColor, "true")
	}

	return args
}
