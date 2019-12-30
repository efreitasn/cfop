package cfp

import (
	"strconv"
	"testing"

	"github.com/efreitasn/customo"
)

func TestBuildOptionOrFlagHelpName(t *testing.T) {
	tests := []struct {
		name        string
		alias       string
		resStyled   string
		resUnstyled string
	}{
		{
			"first",
			"f",
			customo.Format("--first", customo.AttrBold) + ", " + customo.Format("-f", customo.AttrBold),
			"--first, -f",
		},
		{
			"first",
			"f",
			customo.Format("--first", customo.AttrBold) + ", " + customo.Format("-f", customo.AttrBold),
			"--first, -f",
		},
		{
			"first",
			"",
			customo.Format("--first", customo.AttrBold),
			"--first",
		},
		{
			"first",
			"",
			customo.Format("--first", customo.AttrBold),
			"--first",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resStyled, resUnstyled := buildOptionOrFlagHelpName(test.name, test.alias)

			if resStyled != test.resStyled {
				t.Errorf("got %v, want %v", resStyled, test.resStyled)
			}

			if resUnstyled != test.resUnstyled {
				t.Errorf("got %v, want %v", resUnstyled, test.resUnstyled)
			}
		})
	}
}
