package data_test

import (
	"testing"

	"github.com/johnweldon/crawler/data"
)

func TestStack(t *testing.T) {
	d := []string{"this", "is", "a", "test"}
	s := data.StringStack{}
	for _, v := range d {
		s.Push(v)
		t.Logf(s.String())
	}
	for _, _ = range d {
		t.Logf(s.Pop())
	}
	t.Logf("%s", s.Pop())
	t.Logf("%s", s.Pop())
}
