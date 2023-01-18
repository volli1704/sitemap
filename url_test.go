package main_test

import (
	"reflect"
	"testing"

	main "github.com/volli1704/sitemap"
)

var newLnkTests map[string]map[string]main.URL = map[string]map[string]main.URL{
	"is fullfilled url valid": {
		"https://github.com/trololo?a=b&c=d": {
			Protocol: "https",
			Host:     "github.com",
			Path:     "trololo",
			Params:   "a=b&c=d",
		},
	},
	"only path url is valid": {
		"/trololo": {
			Path: "trololo",
		},
	},
	"path with params is valid": {
		"/trololo?a=b&c=d": {
			Path:   "trololo",
			Params: "a=b&c=d",
		},
	},
}

func TestNewLink(t *testing.T) {
	for msg, tst := range newLnkTests {
		for url, lnk := range tst {
			l, _ := main.NewLink(url)
			if !reflect.DeepEqual(l, lnk) {
				t.Errorf("%s: want %+v, got %+v", msg, lnk, l)
			}
		}
	}
}

func TestLinkGetURL(t *testing.T) {
	for msg, tst := range newLnkTests {
		for url := range tst {
			l, _ := main.NewLink(url)
			if l.String() != url {
				t.Errorf("%s: want %s, got %s", msg, url, l)
			}
		}
	}
}
