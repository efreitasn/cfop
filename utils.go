package cfp

import (
	"math"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// getTermNumCols returns the number of columns of the terminal.
func getTermNumCols() (int, error) {
	cols, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, err
	}

	return cols, nil
}

// breakStringIntoPaddedLines breaks str into lines padded with pad padChar whose length <= maxCharsPerLine.
func breakStringIntoPaddedLines(pad int, padChar rune, maxCharsPerLine int, str string) string {
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
