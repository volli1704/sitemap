package url

import (
	"fmt"
	"net/url"
	"strings"
)

// URL is a representation of site URL
// URL can be created from string, has a parent url in context of our url parser and depth
type URL struct {
	Protocol string
	Host     string
	Path     string
	Params   string
	Parent   *URL
	Depth    int
}

// Create URL from string and parent URL
func NewLink(rawURL string, parent *URL) (URL, error) {
	var lnk URL

	prsdURL, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return URL{}, fmt.Errorf("error creating new link from raw string [%s]: %w", rawURL, err)
	}

	lnk.Protocol = prsdURL.Scheme
	lnk.Host = prsdURL.Host
	lnk.Path = prsdURL.Path
	lnk.Params = prsdURL.RawQuery
	lnk.Parent = parent
	if parent != nil {
		lnk.Depth = parent.Depth + 1
	}

	return lnk, nil
}

// Return string representation of url
func (l URL) String() string {
	var url strings.Builder

	if l.Protocol != "" {
		url.WriteString(l.Protocol)
		url.WriteString("://")
	}

	if l.Host != "" {
		url.WriteString(l.Host)
	}

	if len(l.Path) > 0 {
		l.Path = "/" + strings.TrimLeft(l.Path, "/")
		url.WriteString(l.Path)
	}

	if len(l.Params) > 0 {
		url.WriteString("?")
		url.WriteString(l.Params)
	}

	return url.String()
}

// Return string representation with additional info
func (u URL) Log() string {
	var parentURL string
	if u.Parent != nil {
		parentURL = u.Parent.String()
	}

	return fmt.Sprintf("Link: %s\nParent: %s\nDepth: %d\n", u.String(), parentURL, u.Depth)
}

// Unify url to one format in context of one site url parser.
// Example: `/` -> https://parent.com/
// Example: `` -> https://parent.com/some/path?some=query
func (l *URL) Unify(parentURL URL) {
	if l.Path == "./" {
		l.Path = ""
	}

	if l.Host == "" {
		l.Host = parentURL.Host
	}

	if l.Protocol == "" {
		l.Protocol = parentURL.Protocol
	}

	pLen := len(l.Path)
	if pLen > 0 && l.Path[pLen-1:] == "/" {
		l.Path = l.Path[:pLen-1]
	}
}
