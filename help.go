package cfp

import (
	"github.com/efreitasn/customo"
	"io"
	"os"
)

var helpWriter io.Writer = os.Stdout
var helpIndentation = "  "
var numberOfSpacesNameAndDescription = 2

type helper interface {
	help(pp parentParser) string
}

func printHelp(h helper, pp parentParser) {
	helpWriter.Write([]byte(h.help(pp)))
}

func buildOptionOrFlagHelpName(name, alias string) (styled, unstyled string) {
	styled = customo.Format("--"+name, customo.AttrBold)
	unstyled = "--" + name

	if alias != "" {
		styled += ", " + customo.Format("-"+alias, customo.AttrBold)
		unstyled += ", -" + alias
	}

	return styled, unstyled
}
