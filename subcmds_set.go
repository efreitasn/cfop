package cfop

import (
	"fmt"
	"strings"

	"github.com/efreitasn/customo"
)

// Subcmd is a subcmd.
type Subcmd struct {
	Name        string
	Description string
	Parser      Parser
}

// SubcmdsSet is a set of subcmds.
type SubcmdsSet struct {
	items map[string]*Subcmd
}

// NewSubcmdsSet creates a subcmds set.
func NewSubcmdsSet(items ...Subcmd) *SubcmdsSet {
	itemsMap := make(map[string]*Subcmd, len(items))

	for i := range items {
		item := items[i]

		if item.Name == "" {
			panic(ErrMissingSubcmdName)
		}

		if item.Parser == nil {
			panic(ErrMissingSubcmdParser)
		}

		itemsMap[item.Name] = &item
	}

	return &SubcmdsSet{items: itemsMap}
}

// Add adds a subcmd to the set.
// If name == "" or parser == nil, it panics.
func (ss *SubcmdsSet) Add(name, description string, parser Parser) {
	if name == "" {
		panic(ErrMissingSubcmdName)
	}

	if parser == nil {
		panic(ErrMissingSubcmdParser)
	}

	if ss.items == nil {
		ss.items = make(map[string]*Subcmd)
	}

	ss.items[name] = &Subcmd{
		Name:        name,
		Description: description,
		Parser:      parser,
	}
}

// Parse parses a slice of strings.
func (ss *SubcmdsSet) Parse(pp parentParser, strs []string) error {
	if len(strs) == 0 {
		return ErrMissingSubcmd
	}

	str := strs[0]

	if isHelpFlag(str) {
		printHelp(ss, pp)

		return nil
	}

	if isOptionWithValue(str) {
		optName, isAlias := extractOptionName(str)

		return ErrUnexpectedOption{
			OptionName: optName,
			IsAlias:    isAlias,
		}
	}

	if isOptionWithoutValue(str) {
		optName, isAlias := extractOptionName(str)

		return ErrUnexpectedOptionOrFlag{
			OptionOrFlagName: optName,
			IsAlias:          isAlias,
		}
	}

	subcmd, ok := ss.items[str]
	if !ok {
		return ErrUnknownSubcmd{SubcmdName: str}
	}

	return subcmd.Parser.Parse(parentParser{
		parser: ss,
		cmds:   append(pp.cmds, subcmd.Name),
	}, strs[1:])
}

func (ss *SubcmdsSet) help(pp parentParser) string {
	numCols, _ := getTermNumCols()

	sb := strings.Builder{}

	ppDescription := getParentParserDescription(pp)
	if ppDescription != "" {
		sb.WriteString(ppDescription + "\n\n")
	}

	sb.WriteString(fmt.Sprintf("Usage: %v SUBCMD", strings.Join(pp.cmds, " ")))
	sb.WriteString("\n\n")
	sb.WriteString("SUBCMD is one of:\n")

	biggestNameLen := 0

	for _, item := range ss.items {
		if len(item.Name) > biggestNameLen {
			biggestNameLen = len(item.Name)
		}
	}

	for _, item := range ss.items {
		if item.Name == "completion" {
			continue
		}

		nameBold := customo.Format(item.Name, customo.AttrBold)

		sb.WriteString(helpIndentationSpaces + nameBold)
		if item.Description != "" {
			descripFormatted := breakStringIntoPaddedLines(
				helpIndentationNumSpaces+
					numSpacesHelpNameAndDescription+
					biggestNameLen,
				' ',
				numCols,
				item.Description,
			)

			// the new slice was created so that the help name could
			// align with the description.
			sb.Write([]byte(descripFormatted[len(item.Name)+helpIndentationNumSpaces:]))
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}
