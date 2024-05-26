package main

import (
	"log"
	"net/url"
)

func NewURL(s string) url.URL {

	u := url.URL{
		Scheme: "https",
		Host:   "oldschool.runescape.wiki",
		Path:   s,
	}
	log.Printf("url: %s", u.String())
	return u
}
