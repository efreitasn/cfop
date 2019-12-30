package cfp

import (
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var optionWithValueRegExp = regexp.MustCompile("^--?((?:[a-zA-Z]|[0-9])+)=(.+)?$")
var optionWithoutValueRegExp = regexp.MustCompile("^--?((?:[a-zA-Z]|[0-9])+)$")
var optionWithOrWithoutValueRegExp = regexp.MustCompile("^--?((?:[a-zA-Z]|[0-9])+)(?:(?:=(.+))|=)?$")
var helpFlagRegExp = regexp.MustCompile("^(?:--help|-h)$")

func isOptionWithValue(str string) bool {
	return optionWithValueRegExp.MatchString(str)
}

// isOptionWithoutValue returns whether str is a flag, which is an option without a value.
func isOptionWithoutValue(str string) bool {
	return optionWithoutValueRegExp.MatchString(str)
}

func extractOptionName(str string) (name string, isAlias bool) {
	matches := optionWithOrWithoutValueRegExp.FindStringSubmatch(str)

	if len(matches) < 2 {
		return "", false
	}

	name = matches[1]

	return name, !strings.HasPrefix(str, "--")
}

func extractOptionValue(str string) string {
	matches := optionWithValueRegExp.FindStringSubmatch(str)

	if len(matches) < 3 {
		return ""
	}

	return matches[2]
}

// isValidValueForType returns whether value is valid for t and, in case it's true, also the value in the specific type.
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

func isHelpFlag(str string) bool {
	return helpFlagRegExp.MatchString(str)
}

func getTermNumCols() (int, error) {
	cols, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, err
	}

	return cols, nil
}

func breakStringWithPadding(pad, maxCharsPerLine int, padChar rune, str string) string {
	if maxCharsPerLine <= pad {
		return strings.Repeat(string(padChar), maxCharsPerLine)
	}

	if (pad + len(str)) <= maxCharsPerLine {
		return strings.Repeat(string(padChar), pad) + str
	}

	strBs := []byte(str)

	numStrCharsPerLine := maxCharsPerLine - pad
	numStrParts := int(math.Ceil(float64(len(str)) / float64(numStrCharsPerLine)))

	strParts := make([]string, 0, numStrParts)

	for i := 0; i < numStrParts; i++ {
		max := i*numStrCharsPerLine + numStrCharsPerLine

		if max > len(strBs) {
			max = len(strBs)
		}

		partBs := strBs[i*numStrCharsPerLine : max]
		strParts = append(strParts, string(partBs))
	}

	joined := strings.Join(strParts, "\n"+strings.Repeat(string(padChar), pad))

	// remove newline at 0.
	return strings.Repeat(string(padChar), pad) + joined
}
