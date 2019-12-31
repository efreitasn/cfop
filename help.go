package cfop

import (
	"io"
	"os"
	"regexp"

	"github.com/efreitasn/customo"
)

// helpWriter is the writer used to print help messages.
var helpWriter io.Writer = os.Stdout

// helpIndentationNumSpaces is the number of spaces prefixed to some
// lines in a help message.
var helpIndentationNumSpaces = 2

// helpIndentationSpaces is the space character repeated helpIndentationNumSpaces
// times.
var helpIndentationSpaces = "  "

var helpFlagRegExp = regexp.MustCompile("^(?:--help|-h)$")

// numberOfSpacesNameAndDescription is the number of spaces between a
// help name and a description. The helpName could be either an option/
// flag help name or an argument help name.
var numSpacesHelpNameAndDescription = 2

// helper provides a help message to be printed to the user.
type helper interface {
	help(pp parentParser) string
}

// printHelp writes the given helper's help message to helpWriter.
func printHelp(h helper, pp parentParser) {
	helpWriter.Write([]byte(h.help(pp)))
}

// buildOptionOrFlagHelpName buils an option/flag help name given
// a name and an alias.
// It expects to always receive a name != "" and an optional alias.
// It returns both a string with ANSI escape codes and one without
// them.
// For instance, for an option whose name is year and alias is y,
// the help name is: --year, -y.
func buildOptionOrFlagHelpName(name, alias string) (styled, unstyled string) {
	styled = customo.Format("--"+name, customo.AttrBold)
	unstyled = "--" + name

	if alias != "" {
		styled += ", " + customo.Format("-"+alias, customo.AttrBold)
		unstyled += ", -" + alias
	}

	return styled, unstyled
}

// buildArgumentHelpName builds an argument help name given a name.
// It expects to always receive a name != "".
// It returns both a string with ANSI escape codes and one without
// them.
// For instance, for an argument whose name is age, the help name
// is: <age>.
func buildArgumentHelpName(name string) (styled, unstyled string) {
	unstyled = "<" + name + ">"
	styled = customo.Format(unstyled, customo.AttrBold)

	return styled, unstyled
}

func isHelpFlag(str string) bool {
	return helpFlagRegExp.MatchString(str)
}

func findBiggestArgHelpNameLen(args map[string]*CmdArg) int {
	biggest := 0

	for _, arg := range args {
		_, helpName := buildArgumentHelpName(arg.Name)

		if len(helpName) > biggest {
			biggest = len(helpName)
		}
	}

	return biggest
}

func findBiggestOptionOrFlagHelpNameLen(options map[string]*CmdOption, flags map[string]*CmdFlag) int {
	biggest := 0

	for _, option := range options {
		_, helpName := buildOptionOrFlagHelpName(option.Name, option.Alias)

		if len(helpName) > biggest {
			biggest = len(helpName)
		}
	}

	for _, flag := range flags {
		_, helpName := buildOptionOrFlagHelpName(flag.Name, flag.Alias)

		if len(helpName) > biggest {
			biggest = len(helpName)
		}
	}

	return biggest
}
