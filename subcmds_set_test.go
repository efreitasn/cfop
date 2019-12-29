package cfp

import (
	"reflect"
	"testing"
	"time"
)

// mockParser is used only for testing.
type mockParser struct {
	ch chan<- []string
}

func (mp mockParser) Parse(strs []string) error {
	mp.ch <- strs

	return nil
}

func TestSubcmdsSet(t *testing.T) {
	fooCh := make(chan []string, 1)
	fooParser := mockParser{
		ch: fooCh,
	}
	barCh := make(chan []string, 1)
	barParser := mockParser{
		ch: barCh,
	}

	set := NewSubcmdsSet(
		Subcmd{
			Name:   "foo",
			Parser: fooParser,
		},
		Subcmd{
			Name:   "bar",
			Parser: barParser,
		},
	)

	var resStrs []string
	strs := []string{"foo", "--bar", "--foobar=barfoo"}
	expectedStrs := strs[1:]

	set.Parse(strs)

	timer := time.NewTimer(time.Second)

	select {
	case strs := <-fooCh:
		resStrs = strs
	case <-timer.C:
		t.Fatal("timeout")
	}

	if !reflect.DeepEqual(resStrs, expectedStrs) {
		t.Errorf("got %v, want %v", resStrs, expectedStrs)
	}
}
