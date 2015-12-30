package data

import (
	"fmt"
	"net/url"
	"strings"
)

type Link struct {
	name string
	url  *url.URL
	path string
}

func NewLink(name string, url *url.URL, path string) *Link {
	return &Link{name: name, url: url, path: path}
}

func (l *Link) String() string {
	if l == nil {
		return ""
	}
	return fmt.Sprintf("<a href=%q>%s</a>", l.URL(), l.Name())
}

func (l *Link) Name() string {
	if l == nil {
		return ""
	}
	return strings.TrimSpace(l.name)
}

func (l *Link) URL() string {
	if l == nil {
		return ""
	}
	return l.url.String()
}

func (l *Link) Path() string {
	if l == nil {
		return ""
	}
	return strings.ToLower(l.path)
}

type LinkFilter interface {
	Pass(*Link) bool
}

func NewLinkInTableFilter() LinkFilter {
	return &underElementFilter{element: "table"}
}

type underElementFilter struct {
	element string
}

func (f *underElementFilter) Pass(l *Link) bool {
	if f == nil {
		return false
	}
	return strings.Index(l.Path(), "."+f.element) >= 0
}
