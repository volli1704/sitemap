package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	link_parser "github.com/volli1704/link"
)

type LinkMap map[URL]uint16
type LinkQ []URL

func (l LinkMap) String() string {
	var b strings.Builder

	for u, c := range l {
		b.WriteString(u.String())
		b.WriteString(" : ")
		b.WriteString(fmt.Sprintf("%d", c))
		b.WriteRune('\n')
	}

	return b.String()
}

func (q *LinkQ) Dequeue() URL {
	res := (*q)[0]
	*q = (*q)[1:]
	return res
}

var lm LinkMap = make(LinkMap)
var lq LinkQ = make(LinkQ, 0)

func main() {
	lnk := flag.String("url", "https://gobyexample.com/", "Root url of the site from where app will explore link tree")
	depth := flag.Int("depth", 2, "Depth of link research")
	flag.Parse()

	url, err := NewLink(*lnk, 0)
	url.Unify(url)
	if err != nil {
		panic(err)
	}

	lq := append(lq, url)
	lm[url] = 1

	for len(lq) > 0 {
		currLink := lq[0]
		currLink.Unify(url)
		lq.Dequeue()

		resp, err := http.Get(currLink.String())
		if err != nil {
			fmt.Println(err)
			continue
		}

		if currLink.Step >= *depth && *depth != -1 {
			break
		}

		links, err := link_parser.Parse(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, l := range links {
			nl, err := NewLink(l.Href, currLink.Step+1)
			if err != nil {
				fmt.Println(err)
				continue
			}
			nl.Unify(currLink)

			if nl.Host != currLink.Host {
				continue
			}

			if _, ok := lm[nl]; !ok {
				lq = append(lq, nl)
				lm[nl] = 0
			}
			lm[nl]++
		}
	}

	fmt.Println(lm)
}
