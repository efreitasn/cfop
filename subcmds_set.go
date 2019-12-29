package cfp

// Subcmd is a subcmd.
type Subcmd struct {
	Name   string
	Parser Parser
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
func (ss *SubcmdsSet) Parse(strs []string) error {
	if len(strs) == 0 {
		return ErrMissingSubcmd
	}

	str := strs[0]

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

	return subcmd.Parser.Parse(strs[1:])
}
