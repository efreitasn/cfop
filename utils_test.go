package cfop

import (
	"strconv"
	"strings"
	"testing"
)

func TestBreakStringIntoPaddedLines(t *testing.T) {
	tests := []struct {
		pad             int
		maxCharsPerLine int
		str             string
		padChar         rune
		res             string
	}{
		{
			10,
			20,
			"foo bar foo bar",
			' ',
			strings.Repeat(" ", 10) + "foo bar fo\n" + strings.Repeat(" ", 10) + "o bar",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := breakStringIntoPaddedLines(
				test.pad,
				test.padChar,
				test.maxCharsPerLine,
				test.str,
			)

			if res != test.res {
				t.Errorf("got %v, want %v", []byte(res), []byte(test.res))
			}
		})
	}
}
