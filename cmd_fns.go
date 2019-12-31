package cfop

import (
	"regexp"
	"strconv"
	"strings"
)

var optionWithValueRegExp = regexp.MustCompile("^--?((?:[a-zA-Z]|[0-9])+)=(.+)?$")
var optionWithoutValueRegExp = regexp.MustCompile("^--?((?:[a-zA-Z]|[0-9])+)$")
var optionWithOrWithoutValueRegExp = regexp.MustCompile("^--?((?:[a-zA-Z]|[0-9])+)(?:(?:=(.+))|=)?$")

// isValueValidForTermType returns whether value is valid for a given t.
func isValueValidForTermType(t TermType, value string) (interface{}, bool) {
	switch t {
	case TermString:
		return value, true
	case TermInt:
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, false
		}

		return n, true
	case TermFloat:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, false
		}

		return n, true
	default:
		panic(ErrInvalidTermType)
	}
}

// isOptionWithValue returns whether str is a option with value, e.g. --name=John
func isOptionWithValue(str string) bool {
	return optionWithValueRegExp.MatchString(str)
}

// isOptionWithoutValue returns whether str is a flag, which is an option without a value, e.g. --lines.
func isOptionWithoutValue(str string) bool {
	return optionWithoutValueRegExp.MatchString(str)
}

// extractOptionName extracts the name from an option, e.g. year out of --year=1990.
func extractOptionName(str string) (name string, isAlias bool) {
	matches := optionWithOrWithoutValueRegExp.FindStringSubmatch(str)

	if len(matches) < 2 {
		return "", false
	}

	name = matches[1]

	return name, !strings.HasPrefix(str, "--")
}

// extractOptionName extracts the value from an option, e.g. 1990 out of --year=1990.
func extractOptionValue(str string) string {
	matches := optionWithValueRegExp.FindStringSubmatch(str)

	if len(matches) < 3 {
		return ""
	}

	return matches[2]
}
