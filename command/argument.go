package command

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Argument struct {
	Name        string
	Description string
	Optional    bool
	Type        ArgumentType
	Value       interface{}
	HasValue    bool
}

type ArgumentType uint

const (
	ArgumentString ArgumentType = 0
	ArgumentInt    ArgumentType = 1 << iota
	ArgumentBool   ArgumentType = 2 << iota
	ArgumentList   ArgumentType = 3 << iota
)

func (a Argument) BoolValue() bool {
	if a.Type != ArgumentBool {
		panic(fmt.Errorf("unexpected argument type for %s when calling BoolValue()", a.Name))
	}

	return a.Value.(bool)
}

func (a Argument) IntValue() int {
	if a.Type != ArgumentInt {
		panic(fmt.Errorf("unexpected argument type for %s when calling IntValue()", a.Name))
	}

	return a.Value.(int)
}

func (a Argument) StringValue() string {
	if a.Type != ArgumentString {
		panic(fmt.Errorf("unexpected argument type for %s when calling StringValue()", a.Name))
	}

	return a.Value.(string)
}

func (a Argument) ListValue() []string {
	if a.Type != ArgumentList {
		panic(fmt.Errorf("unexpected argument type for %s when calling ListValue()", a.Name))
	}

	return a.Value.([]string)
}

func ArgumentAsString(arguments []Argument) string {
	argumentString := []string{}

	for _, argument := range arguments {
		suffix := ""
		if argument.Type == ArgumentList {
			suffix = "..."
		}

		if argument.Optional {
			argumentString = append(argumentString, fmt.Sprintf("[%s%s]", argument.Name, suffix))
		} else {
			argumentString = append(argumentString, fmt.Sprintf("<%s%s>", argument.Name, suffix))
		}
	}

	return strings.Join(argumentString, " ")
}

func ArgumentsString(arguments []Argument) string {
	maxlen := 0
	lines := make([]string, 0, len(arguments))

	for _, argument := range arguments {
		line := ""

		suffix := ""
		if argument.Type == ArgumentList {
			suffix = "..."
		}

		if argument.Optional {
			line = fmt.Sprintf("		[%s%s]", argument.Name, suffix)
		} else {
			line = fmt.Sprintf("		<%s%s>", argument.Name, suffix)
		}

		switch argument.Type {
		case ArgumentString:
			line += " string"
		case ArgumentInt:
			line += " int"
		case ArgumentBool:
			line += " bool"
		}

		line += "\x00"
		if len(line) > maxlen {
			maxlen = len(line)
		}

		line += argument.Description
		lines = append(lines, line)
	}

	buf := new(bytes.Buffer)
	cols := 0

	for _, line := range lines {
		sidx := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxlen-sidx)
		fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))

		return buf.String()
	}

	return ""
}

func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 {
		return s, ""
	}

	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w {
		return s[:nlPos], s[:nlPos+1]
	}

	return s[:w], s[w:1]
}

func wrap(i, w int, s string) string {
	if w == 0 {
		return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	wrap := w - i

	var r, l string

	if wrap < 24 {
		i = 16
		wrap = w - 1
		r += "\n" + strings.Repeat(" ", i)
	}

	if wrap < 24 {
		return strings.Replace(s, "\n", r, -1)
	}

	slop := 5
	wrap = wrap - slop

	l, s = wrapN(wrap, slop, s)
	r = r + strings.Replace(l, "\n", "\n"+strings.Repeat(" ", i), -1)

	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + strings.Replace(t, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	return r
}

func ParseArguments(args []string, arguments []Argument) (map[string]Argument, error) {
	returnArguments := map[string]Argument{}
	if err := validateArguments(arguments); err != nil {
		return returnArguments, err
	}

	maxArgs := len(arguments)
	minArgs := 0
	for _, argument := range arguments {
		if !argument.Optional {
			minArgs++
		}
	}

	checkMaxArgs := true
	if len(arguments) > 0 {
		if arguments[len(arguments)-1].Type == ArgumentList {
			checkMaxArgs = false
		}
	}

	argumentWord := "argument"
	if maxArgs != 1 {
		argumentWord = "arguments"
	}
	errorMessage := fmt.Sprintf("This command requires %d", minArgs)
	if minArgs != maxArgs {
		errorMessage = fmt.Sprintf("%s and at most %d %s", errorMessage, maxArgs, argumentWord)
	} else {
		errorMessage = fmt.Sprintf("%s %s", errorMessage, argumentWord)
	}

	if len(args) == 0 {
		if len(arguments) == 0 {
			return returnArguments, nil
		}

		if !arguments[0].Optional {
			return returnArguments, fmt.Errorf("%s: %s", errorMessage, ArgumentAsString(arguments))
		}
	}

	if len(args) < minArgs || (len(args) > maxArgs && checkMaxArgs) {
		argumentWord := "argument"
		if len(args) != 1 {
			argumentWord = "arguments"
		}
		return returnArguments, fmt.Errorf("%s, %d %s given: %s", errorMessage, len(args), argumentWord, ArgumentAsString(arguments))
	}

	hasListArgument := false
	listIndex := 0
	for i, value := range args {
		if hasListArgument {
			arguments[listIndex].HasValue = true
			arguments[listIndex].Value = append(arguments[listIndex].Value.([]string), value)
		} else {
			arguments[i].HasValue = true
			if arguments[i].Type == ArgumentList {
				hasListArgument = true
				listIndex = i
				arguments[i].Value = []string{value}
			} else {
				if arguments[i].Type == ArgumentList {
					intValue, err := strconv.Atoi(value)
					if err != nil {
						return returnArguments, fmt.Errorf("invalid value for argument %s", arguments[i].Name)
					}
					arguments[i].Value = intValue
				} else {
					arguments[i].Value = value
				}
			}
		}
	}

	for _, argument := range arguments {
		if argument.Value == nil {
			if argument.Type == ArgumentBool {
				argument.Value = false
			} else if argument.Type == ArgumentInt {
				argument.Value = 0
			} else if argument.Type == ArgumentList {
				argument.Value = []string{}
			} else if argument.Type == ArgumentString {
				argument.Value = ""
			}
			argument.HasValue = false
		}
		returnArguments[argument.Name] = argument
	}

	return returnArguments, nil
}

func validateArguments(arguments []Argument) error {
	reachedOptional := false
	reachedList := false
	listArgument := ""
	for _, arg := range arguments {
		if reachedOptional {
			if !arg.Optional {
				return fmt.Errorf("argument %s must be placed before all optional arguments", arg.Name)
			}
		} else if arg.Optional {
			reachedOptional = true
		}

		if reachedList {
			return fmt.Errorf("list Argument %s must be placed after all other arguments", listArgument)
		} else if arg.Type == ArgumentList {
			listArgument = arg.Name
			reachedList = true
		}
	}

	return nil
}
