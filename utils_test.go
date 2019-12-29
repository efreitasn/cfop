package cfp

import (
	"strconv"
	"testing"
)

func TestIsOptionWithValue(t *testing.T) {
	tests := []struct {
		str string
		res bool
	}{
		{"--opt=300", true},
		{"--opt=", true},
		{"--opt", false},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := isOptionWithValue(test.str)

			if res != test.res {
				t.Errorf("got %v, want %v", res, test.res)
			}
		})
	}
}

func TestIsOptionWithoutValue(t *testing.T) {
	tests := []struct {
		str string
		res bool
	}{
		{"--opt=300", false},
		{"--opt", true},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := isOptionWithoutValue(test.str)

			if res != test.res {
				t.Errorf("got %v, want %v", res, test.res)
			}
		})
	}
}

func TestExtractOptionName(t *testing.T) {
	tests := []struct {
		str     string
		res     string
		isAlias bool
	}{
		{"--opt=300", "opt", false},
		{"--opt=", "opt", false},
		{"--opt", "opt", false},
		{"-o", "o", true},
		{"-o=saas", "o", true},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, isAlias := extractOptionName(test.str)

			if res != test.res {
				t.Errorf("got %v, want %v", res, test.res)
			}

			if isAlias != test.isAlias {
				t.Errorf("got %v, want %v", isAlias, test.isAlias)
			}
		})
	}
}

func TestExtractOptionValue(t *testing.T) {
	tests := []struct {
		str string
		res string
	}{
		{"--opt=300", "300"},
		{"--opt", ""},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := extractOptionValue(test.str)

			if res != test.res {
				t.Errorf("got %v, want %v", res, test.res)
			}
		})
	}
}

func TestIsValueValidForTermType(t *testing.T) {
	tests := []struct {
		t        TermType
		val      string
		resVal   interface{}
		resValid bool
	}{
		{
			TermInt,
			"20",
			20,
			true,
		},
		{
			TermInt,
			"sakosa",
			nil,
			false,
		},
		{
			TermString,
			"sakosa",
			"sakosa",
			true,
		},
		{
			TermFloat,
			"20.50",
			20.50,
			true,
		},
		{
			TermFloat,
			"ksakosaoksa",
			nil,
			false,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			val, valid := isValueValidForTermType(test.t, test.val)

			if valid != test.resValid {
				t.Errorf("got %v, want %v", valid, test.resValid)
			}

			if val != test.resVal {
				t.Errorf("got %v, want %v", val, test.resVal)
			}
		})
	}
}
