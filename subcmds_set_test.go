package cfop

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

// mockParser is used only for testing.
type mockParser struct {
	ch chan<- []string
}

func (mp mockParser) Parse(pp parentParser, strs []string) error {
	mp.ch <- strs

	return nil
}

func TestSubcmdsSet(t *testing.T) {
	tests := []struct {
		subcmdsNames []string
		strs         []string
		err          error
	}{
		{
			[]string{"foo", "bar"},
			[]string{"foo", "--bar", "--foobar=barfoo"},
			nil,
		},
		{
			[]string{"foo", "bar"},
			[]string{},
			ErrMissingSubcmd,
		},
		{
			[]string{"foo", "bar"},
			[]string{"--some"},
			ErrUnexpectedOptionOrFlag{
				OptionOrFlagName: "some",
				IsAlias:          false,
			},
		},
		{
			[]string{"foo", "bar"},
			[]string{"--some=that"},
			ErrUnexpectedOption{
				OptionName: "some",
				IsAlias:    false,
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			subcmdsChs := make(map[string]chan []string, len(test.subcmdsNames))
			subcmdsParsers := make(map[string]Parser, len(test.subcmdsNames))
			set := NewSubcmdsSet()

			for _, subcmdName := range test.subcmdsNames {
				subcmdsChs[subcmdName] = make(chan []string, 1)
				subcmdsParsers[subcmdName] = mockParser{
					ch: subcmdsChs[subcmdName],
				}
				set.Add(subcmdName, "", subcmdsParsers[subcmdName])
			}

			err := set.Parse(parentParser{
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

			var resStrs []string
			expectedStrs := test.strs[1:]

			timer := time.NewTimer(time.Second)

			select {
			case strs := <-subcmdsChs[test.strs[0]]:
				resStrs = strs
			case <-timer.C:
				t.Fatal("timeout")
			}

			if !reflect.DeepEqual(resStrs, expectedStrs) {
				t.Errorf("got %v, want %v", resStrs, expectedStrs)
			}
		})
	}
}
