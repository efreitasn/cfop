package cfop

import (
	"fmt"
	"sort"
	"strings"
)

// parentParser is a reference to the previous parser.
type parentParser struct {
	// cmds is a slice containing the name of each cmd executed thus far.
	cmds   []string
	parser Parser
}

// Parser parses a slice of strings.
type Parser interface {
	Parse(pp parentParser, strs []string) error
}

// rootCmd is a parser representing the root cmd.
type rootCmd struct {
	name        string
	description string
}

func (rp *rootCmd) Parse(pp parentParser, strs []string) error {
	return nil
}

// getParentParserDescription returns the description of p.
// If p doesn't have a description, an empty string is returned.
func getParentParserDescription(pp parentParser) string {
	switch p := pp.parser.(type) {
	case *SubcmdsSet:
		if item := p.items[pp.cmds[len(pp.cmds)-1]]; item.Description != "" {
			return item.Description
		}
	case *rootCmd:
		if p.description != "" {
			return p.description
		}
	}

	return ""
}

func introspectParser(strs []string, p Parser) []string {
	res := []string{"--help", "-h"}

	for i := 0; i < len(strs); i++ {
		str := strs[i]

		if str == "" {
			break
		}

		switch cmdOrSet := p.(type) {
		case *Cmd:
			continue
		case *SubcmdsSet:
			item, ok := cmdOrSet.items[str]
			if !ok {
				return res
			}

			p = item.Parser
			continue
		}
	}

	switch cmdOrSet := p.(type) {
	case *Cmd:
		for _, opt := range cmdOrSet.options {
			res = append(res, "--"+opt.Name)

			if opt.Alias != "" {
				res = append(res, "-"+opt.Alias)
			}
		}
		for _, flag := range cmdOrSet.flags {
			res = append(res, "--"+flag.Name)

			if flag.Alias != "" {
				res = append(res, "-"+flag.Alias)
			}
		}
	case *SubcmdsSet:
		for _, item := range cmdOrSet.items {
			res = append(res, item.Name)
		}
	}

	sort.Strings(res)

	return res
}

// Init initiates the parsing of the CLI.
// The first item in strs is ignored, so that
// os.Args can be used as the strs' value, which
// is generally the case.
func Init(name, description string, strs []string, p Parser) error {
	if name == "" {
		return ErrMissingRootCmdName
	}

	rp := &rootCmd{
		name:        name,
		description: description,
	}

	var newStrs []string

	if len(strs) == 1 {
		newStrs = make([]string, 0)
	} else {
		newStrs = strs[1:]
	}

	// Introspection
	if len(newStrs) > 0 && newStrs[0] == "__introspect__" {
		fmt.Println(strings.Join(
			introspectParser(newStrs[1:], p),
			" ",
		))

		return nil
	}

	return p.Parse(parentParser{
		cmds:   []string{name},
		parser: rp,
	}, newStrs)
}
