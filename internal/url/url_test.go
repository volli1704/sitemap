package url_test

import (
	"reflect"
	"testing"

	url "github.com/volli1704/sitemap/internal/url"
)

var newLnkTests map[string]map[string]url.URL = map[string]map[string]url.URL{
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
		for u, lnk := range tst {
			l, _ := url.NewLink(u, nil)
			if !reflect.DeepEqual(l, lnk) {
				t.Errorf("%s: want %+v, got %+v", msg, lnk, l)
			}
		}
	}
}

func TestLinkGetURL(t *testing.T) {
	for msg, tst := range newLnkTests {
		for u := range tst {
			l, _ := url.NewLink(u, nil)
			if l.String() != u {
				t.Errorf("%s: want %s, got %s", msg, u, l)
			}
		}
	}
}
