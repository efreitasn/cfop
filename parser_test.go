package cfp

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGetParentParserDescription(t *testing.T) {
	tests := []struct {
		pp  parentParser
		res string
	}{
		{
			parentParser{
				cmds: []string{"foo", "bar"},
				parser: NewSubcmdsSet(
					Subcmd{
						Name:        "bar",
						Description: "some stuff",
						Parser: NewCmd(CmdConfig{
							Fn: func(cts *CmdTermsSet) {},
						}),
					},
				),
			},
			"some stuff",
		},
	}

	for _, test := range tests {
		res := getParentParserDescription(test.pp)

		if res != test.res {
			t.Errorf("got %v, want %v", res, test.res)
		}
	}
}

func TestInstrospectParser(t *testing.T) {
	tests := []struct {
		p    Parser
		strs []string
		res  []string
	}{
		{
			NewSubcmdsSet(
				Subcmd{
					Name: "foo",
					Parser: NewCmd(CmdConfig{
						Fn: func(cts *CmdTermsSet) {},
						Flags: []CmdFlag{
							CmdFlag{
								Name:  "year",
								Alias: "y",
							},
						},
					}),
				},
			),
			[]string{"foo", "--y"},
			[]string{"--year", "-y"},
		},
		{
			NewSubcmdsSet(
				Subcmd{
					Name: "foo",
					Parser: NewCmd(CmdConfig{
						Fn: func(cts *CmdTermsSet) {},
						Options: []CmdOption{
							CmdOption{
								Name:  "age",
								Alias: "a",
								T:     TermInt,
							},
						},
						Flags: []CmdFlag{
							CmdFlag{
								Name:  "year",
								Alias: "y",
							},
						},
					}),
				},
			),
			[]string{"foo"},
			[]string{"--age", "--year", "-a", "-y"},
		},
		{
			NewSubcmdsSet(
				Subcmd{
					Name: "foo",
					Parser: NewCmd(CmdConfig{
						Fn: func(cts *CmdTermsSet) {},
						Options: []CmdOption{
							CmdOption{
								Name:  "age",
								Alias: "a",
								T:     TermInt,
							},
						},
						Flags: []CmdFlag{
							CmdFlag{
								Name:  "year",
								Alias: "y",
							},
						},
					}),
				},
			),
			[]string{"foo", "--age", "1990"},
			[]string{"--age", "--year", "-a", "-y"},
		},
		{
			NewSubcmdsSet(
				Subcmd{
					Name: "foo",
					Parser: NewCmd(CmdConfig{
						Fn: func(cts *CmdTermsSet) {},
						Options: []CmdOption{
							CmdOption{
								Name:  "age",
								Alias: "a",
								T:     TermInt,
							},
						},
						Flags: []CmdFlag{
							CmdFlag{
								Name:  "year",
								Alias: "y",
							},
						},
					}),
				},
			),
			[]string{},
			[]string{"foo"},
		},
		{
			NewSubcmdsSet(
				Subcmd{
					Name: "foo",
					Parser: NewCmd(CmdConfig{
						Fn: func(cts *CmdTermsSet) {},
						Options: []CmdOption{
							CmdOption{
								Name:  "age",
								Alias: "a",
								T:     TermInt,
							},
						},
						Flags: []CmdFlag{
							CmdFlag{
								Name:  "year",
								Alias: "y",
							},
						},
					}),
				},
			),
			[]string{""},
			[]string{"foo"},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := introspectParser(test.strs, test.p)

			if !reflect.DeepEqual(res, test.res) {
				t.Errorf("got %v, want %v", res, test.res)
			}
		})
	}
}
