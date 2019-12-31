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
func TestBuildArgumentHelpName(t *testing.T) {
	tests := []struct {
		name        string
		resStyled   string
		resUnstyled string
	}{
		{
			"first",
			customo.Format("<first>", customo.AttrBold),
			"<first>",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resStyled, resUnstyled := buildArgumentHelpName(test.name)

			if resStyled != test.resStyled {
				t.Errorf("got %v, want %v", resStyled, test.resStyled)
			}

			if resUnstyled != test.resUnstyled {
				t.Errorf("got %v, want %v", resUnstyled, test.resUnstyled)
			}
		})
	}
}

func TestIsHelpFlag(t *testing.T) {
	tests := []struct {
		str string
		res bool
	}{
		{"--foo", false},
		{"--help", true},
		{"-h", true},
		{"-help", false},
		{"--h", false},
		{"some", false},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := isHelpFlag(test.str)

			if res != test.res {
				t.Errorf("got %v, want %v", res, test.res)
			}
		})
	}
}

func TestFindBiggestOptionOrFlagHelpNameLen(t *testing.T) {
	tests := []struct {
		options map[string]*CmdOption
		flags   map[string]*CmdFlag
		res     int
	}{
		{
			map[string]*CmdOption{
				"year": &CmdOption{
					Name:  "year",
					Alias: "y",
				},
				"age": &CmdOption{
					Name:  "age",
					Alias: "",
				},
			},
			map[string]*CmdFlag{
				"lines": &CmdFlag{
					Name:  "lines",
					Alias: "l",
				},
			},
			len("--lines, -l"),
		},
		{
			map[string]*CmdOption{},
			map[string]*CmdFlag{},
			0,
		},
		{
			nil,
			nil,
			0,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := findBiggestOptionOrFlagHelpNameLen(test.options, test.flags)

			if res != test.res {
				t.Errorf("got %v, want %v", res, test.res)
			}
		})
	}
}

func TestFindBiggestArgHelpNameLen(t *testing.T) {
	tests := []struct {
		args map[string]*CmdArg
		res  int
	}{
		{
			map[string]*CmdArg{
				"year": &CmdArg{
					Name: "year",
				},
				"foobarbarfoo": &CmdArg{
					Name: "foobarbarfoo",
				},
				"foobar": &CmdArg{
					Name: "foobar",
				},
			},
			len("<foobarbarfoo>"),
		},
		{
			map[string]*CmdArg{},
			0,
		},
		{
			nil,
			0,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := findBiggestArgHelpNameLen(test.args)

			if res != test.res {
				t.Errorf("got %v, want %v", res, test.res)
			}
		})
	}
}
