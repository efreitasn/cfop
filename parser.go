package cfp

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

// rootCmd is a parser representing the first cmd of the chain.
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

	return p.Parse(parentParser{
		cmds:   []string{name},
		parser: rp,
	}, newStrs)
}
