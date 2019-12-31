package cfp

import "testing"

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
