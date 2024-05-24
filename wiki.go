package main

import "strings"

type URL struct {
	base  string
	path  string
	full  string
	valid bool
}

var InvalidChars = []string{
	",", ".", "\"", "{", "}", "[", "]",
}

func NewURL(base string, path string) *URL {
	if base == "" {
		base = "https://oldschool.runescape.wiki/"
	}
	return &URL{path: path, base: base, full: base + path}
}

func (u *URL) isValid() {
	for _, char := range InvalidChars {
		if strings.Contains(u.path, char) {
			u.valid = false
			break
		}
	}
	u.valid = true

}
