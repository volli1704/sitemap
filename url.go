package main

import (
	"fmt"
	"net/url"
	"strings"
)

type URL struct {
	Protocol string
	Host     string
	Path     string
	Params   string
	Step     int
}

func NewLink(rawURL string, params ...int) (URL, error) {
	var lnk URL

	prsdURL, err := url.Parse(rawURL)
	if err != nil {
		return URL{}, fmt.Errorf("error creating new link from raw string [%s]: %w", rawURL, err)
	}

	lnk.Protocol = prsdURL.Scheme
	lnk.Host = prsdURL.Host
	lnk.Path = prsdURL.Path
	lnk.Params = prsdURL.RawQuery
	if len(params) > 0 {
		lnk.Step = params[0]
	}

	return lnk, nil
}

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
