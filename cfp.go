package cfp

// Parser parses a slice of strings.
type Parser interface {
	Parse([]string) error
}
