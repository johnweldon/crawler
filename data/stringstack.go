package data

import "strings"

type StringStack struct {
	data []string
}

func (s *StringStack) Push(v string) { s.data = append(s.data, v) }

func (s *StringStack) Pop() string {
	if len(s.data) == 0 {
		return ""
	}
	x, d := s.data[len(s.data)-1], s.data[:len(s.data)-1]
	s.data = d
	return x
}

func (s *StringStack) Peek() string {
	if len(s.data) == 0 {
		return ""
	}
	return s.data[len(s.data)-1]
}

func (s StringStack) String() string {
	if len(s.data) == 0 {
		return ""
	}
	return strings.Join(s.data, ".")
}
