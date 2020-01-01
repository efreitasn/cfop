package cfop

import (
	"strconv"
	"testing"
	"time"
)

func TestCmd(t *testing.T) {
	tests := []struct {
		config     CmdConfig
		strs       []string
		err        error
		intOpts    map[string]int
		floatOpts  map[string]float64
		stringOpts map[string]string
		intArgs    map[string]int
		floatArgs  map[string]float64
		stringArgs map[string]string
		flags      map[string]bool
	}{
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"name", "n", "", TermString, false},
					{"year", "y", "", TermInt, true},
					{"salary", "sl", "", TermFloat, true},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
				Args: []CmdArg{
					{"first", "", TermString},
				},
			},
			strs: []string{"-n", "John", "-y=1990", "foobar", "-l", "--salary", "500.85"},
			err:  nil,
			stringOpts: map[string]string{
				"name": "John",
			},
			intOpts: map[string]int{
				"year": 1990,
			},
			floatOpts: map[string]float64{
				"salary": 500.85,
			},
			flags: map[string]bool{
				"line": true,
			},
			stringArgs: map[string]string{"first": "foobar"},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"title", "t", "", TermString, false},
				},
				Args: []CmdArg{
					{"first", "", TermFloat},
				},
			},
			strs: []string{"--title", "salary", "180.87"},
			err:  nil,
			stringOpts: map[string]string{
				"title": "salary",
			},
			floatArgs: map[string]float64{"first": 180.87},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"title", "t", "", TermString, false},
				},
				Args: []CmdArg{
					{"first", "", TermInt},
				},
			},
			strs: []string{"--title", "salary", "180.87"},
			err: ErrArgumentExpectsDifferentValueType{
				ArgumentPos:  0,
				ArgumentName: "first",
				Value:        "180.87",
				ExpectedType: TermInt,
			},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"title", "t", "", TermString, false},
				},
				Args: []CmdArg{
					{"first", "", TermInt},
				},
			},
			strs: []string{"--title", "first", "985"},
			err:  nil,
			stringOpts: map[string]string{
				"title": "first",
			},
			intArgs: map[string]int{"first": 985},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"name", "n", "", TermString, false},
					{"year", "y", "", TermInt, true},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
				Args: []CmdArg{
					{"first", "", TermString},
				},
			},
			strs: []string{"-n", "John", "--year=foo", "foobar"},
			err: ErrOptionExpectsDifferentValueType{
				OptionName:   "year",
				ExpectedType: TermInt,
			},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"name", "n", "", TermString, false},
					{"year", "y", "", TermInt, true},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
				Args: []CmdArg{
					{"first", "", TermString},
				},
			},
			strs: []string{"-n", "John", "--year", "foo", "foobar"},
			err: ErrOptionExpectsDifferentValueType{
				OptionName:   "year",
				ExpectedType: TermInt,
			},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"name", "n", "", TermString, false},
					{"year", "y", "", TermInt, true},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
			},
			strs: []string{"-n", "John", "--year"},
			err: ErrOptionsExpectsAValue{
				OptionName: "year",
			},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"name", "n", "", TermString, false},
					{"year", "y", "", TermInt, true},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
			},
			strs: []string{"-n", "John", "--year="},
			err: ErrOptionsExpectsAValue{
				OptionName: "year",
			},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"name", "n", "", TermString, false},
					{"year", "y", "", TermInt, true},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
				Args: []CmdArg{
					{"first", "", TermString},
				},
			},
			strs: []string{"-n", "John", "--year=1990", "--what", "foobar"},
			err:  ErrUnexpectedOptionOrFlag{OptionOrFlagName: "what"},
		},
		{
			config: CmdConfig{
				Options: []CmdOption{
					{"name", "n", "", TermString, false},
					{"year", "y", "", TermInt, true},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
			},
			strs: []string{"-n", "John", "--year=1990", "foobar"},
			err:  ErrUnexpectedArgument{Argument: "foobar"},
		},
		{
			config: CmdConfig{
				Args: []CmdArg{
					{"first", "", TermFloat},
					{"Second", "", TermInt},
				},
				Flags: []CmdFlag{
					{"line", "l", ""},
				},
			},
			strs: []string{"-l"},
			err:  ErrMissingArguments,
		},
		{
			config: CmdConfig{
				Args: []CmdArg{
					{"first", "", TermInt},
				},
				Options: []CmdOption{
					{"name", "n", "", TermString, true},
				},
			},
			strs: []string{"20"},
			err:  ErrRequiredOptionNotProvided{OptionName: "name"},
		},
		{
			config: CmdConfig{},
			strs:   []string{"--name", "John"},
			err:    ErrUnexpectedOptionOrFlag{OptionOrFlagName: "name"},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ch := make(chan *CmdTermsSet, 1)
			newConfig := test.config

			newConfig.Fn = func(set *CmdTermsSet) {
				ch <- set
			}

			cmd := NewCmd(newConfig)

			err := cmd.Parse(parentParser{
				parser: &rootCmd{
					name: "testing",
				},
				cmds: []string{"testing"},
			}, test.strs)

			if err != test.err {
				t.Fatalf("got %v, want %v", err, test.err)
			}

			if err != nil {
				return
			}

			var set *CmdTermsSet
			timer := time.NewTimer(time.Second)

			select {
			case s := <-ch:
				set = s
			case <-timer.C:
				t.Fatal("timeout")
			}

			if test.intOpts != nil {
				for name, expectedRes := range test.intOpts {
					res := set.GetOptInt(name)

					if res != expectedRes {
						t.Fatalf("got %v, want %v", res, expectedRes)
					}
				}
			}

			if test.floatOpts != nil {
				for name, expectedRes := range test.floatOpts {
					res := set.GetOptFloat(name)

					if res != expectedRes {
						t.Fatalf("got %v, want %v", res, expectedRes)
					}
				}
			}

			if test.stringOpts != nil {
				for name, expectedRes := range test.stringOpts {
					res := set.GetOptString(name)

					if res != expectedRes {
						t.Fatalf("got %v, want %v", res, expectedRes)
					}
				}
			}

			if test.intArgs != nil {
				for name, expectedRes := range test.intArgs {
					res := set.GetArgInt(name)

					if res != expectedRes {
						t.Fatalf("got %v, want %v", res, expectedRes)
					}
				}
			}

			if test.floatArgs != nil {
				for name, expectedRes := range test.floatArgs {
					res := set.GetArgFloat(name)

					if res != expectedRes {
						t.Fatalf("got %v, want %v", res, expectedRes)
					}
				}
			}

			if test.stringArgs != nil {
				for name, expectedRes := range test.stringArgs {
					res := set.GetArgString(name)

					if res != expectedRes {
						t.Fatalf("got %v, want %v", res, expectedRes)
					}
				}
			}

			if test.flags != nil {
				for name, expectedRes := range test.flags {
					res := set.GetFlag(name)

					if res != expectedRes {
						t.Fatalf("got %v, want %v", res, expectedRes)
					}
				}
			}
		})
	}
}
