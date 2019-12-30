package cfp

// parentParser holds a reference to the previous parser.
type parentParser struct {
	// cmds is a slice containing the name of each cmd executed thus far.
	cmds   []string
	parser Parser
}

// Parser parses a slice of strings.
type Parser interface {
	Parse(pp parentParser, strs []string) error
}

type rootCmd struct {
	name        string
	description string
}

func (rp *rootCmd) Parse(pp parentParser, strs []string) error {
	return nil
}

// ParseRoot initiates the parsing.
func ParseRoot(name, description string, strs []string, p Parser) error {
	rp := &rootCmd{
		name:        name,
		description: description,
	}

	return p.Parse(parentParser{
		cmds:   []string{name},
		parser: rp,
	}, strs)
}
