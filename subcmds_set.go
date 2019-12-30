package cfp

import (
	"fmt"
	"github.com/efreitasn/customo"
	"strings"
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

		itemsMap[item.Name] = &item
	}

	return &SubcmdsSet{items: itemsMap}
}

// Add adds a subcmd to the set.
func (ss *SubcmdsSet) Add(item Subcmd) {
	if ss.items == nil {
		ss.items = make(map[string]*Subcmd)
	}

	ss.items[item.Name] = &item
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
	sb := strings.Builder{}

	// Parent cmd's description.
	switch p := pp.parser.(type) {
	case *SubcmdsSet:
		if item := p.items[pp.cmds[len(pp.cmds)-1]]; item.Description != "" {
			sb.WriteString(item.Description + "\n\n")
		}
	case *rootCmd:
		if p.description != "" {
			sb.WriteString(p.description + "\n\n")
		}
	}

	sb.WriteString(fmt.Sprintf("Usage: %v SUBCMD", strings.Join(pp.cmds, " ")))
	sb.WriteString("\n\n")
	sb.WriteString("SUBCMD is one of:\n")

	largestNameLen := 0

	for _, item := range ss.items {
		if len(item.Name) > largestNameLen {
			largestNameLen = len(item.Name)
		}
	}

	for _, item := range ss.items {
		nameBold := customo.Format(item.Name, customo.AttrBold)

		sb.WriteString(helpIndentation + nameBold)
		if item.Description != "" {
			spaces := strings.Repeat(" ", 5+largestNameLen-len(item.Name))

			sb.WriteString(spaces + item.Description)
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
