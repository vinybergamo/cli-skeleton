package command

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

func MergeAutocompleteFlags(flags ...complete.Flags) complete.Flags {
	merged := make(map[string]complete.Predictor, len(flags))
	for _, f := range flags {
		for k, v := range f {
			merged[k] = v
		}
	}
	return merged
}

func CommandErrorText(cmd NamedCommand) string {
	appName := os.Getenv("CLI_APP_NAME")
	return fmt.Sprintf("For additional help try '%s %s --help'", appName, cmd.Name())
}

type uiErrorWriter struct {
	ui  cli.Ui
	buf bytes.Buffer
}

func (w *uiErrorWriter) Write(data []byte) (int, error) {
	read := 0
	for len(data) != 0 {
		a, token, err := bufio.ScanLines(data, false)
		if err != nil {
			return read, err
		}

		if a == 0 {
			r, err := w.buf.Write(data)
			return read + r, err
		}

		w.ui.Error(w.buf.String() + string(token))
		data = data[a:]
		w.buf.Reset()
		read += a
	}

	return read, nil
}

func (w *uiErrorWriter) Close() error {
	if w.buf.Len() != 0 {
		w.ui.Error(w.buf.String())
		w.buf.Reset()
	}

	return nil
}
