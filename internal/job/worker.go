package job

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	link_parser "github.com/volli1704/link"
	url "github.com/volli1704/sitemap/internal/url"
)

const HTTP_CLIENT_TIMEOUT = 3 * time.Second

// Worker parse url from InChan and pass all url from page to OutChan.
// Also it utilize Sync chan to push value there after start/end of processing URL
// so runner know when all workers in idle and stops `url loop`
type worker struct {
	id       int8
	SyncChan chan bool
	InChan   chan url.URL
	OutChan  chan url.URL
}

func NewWorker(id int8, inC, outC chan url.URL, syncC chan bool) worker {
	return worker{
		id,
		syncC,
		inC,
		outC,
	}
}

func (j *worker) run() {
	fmt.Printf("Worker %d started\n", j.id)

	for {
		lnk := <-j.InChan
		j.SyncChan <- false

		fmt.Printf("Processing link %s\n", lnk)

		client := &http.Client{
			Timeout: HTTP_CLIENT_TIMEOUT,
		}
		resp, err := client.Get(lnk.String())
		if err != nil {
			fmt.Printf("%s url get error: %s\n", lnk, err)
			<-j.SyncChan
			continue
		}

		suitContent := false
		for _, h := range resp.Header["Content-Type"] {
			if strings.HasPrefix(h, "text/html") {
				suitContent = true
			}
		}

		if resp == nil || !suitContent {
			<-j.SyncChan
			continue
		}
		links, err := link_parser.Parse(resp.Body)
		if err != nil {
			<-j.SyncChan
			fmt.Printf("%s page parse error: %s\n", lnk, err)
			fmt.Println("parent url:", lnk.Log())
			continue
		}

		for _, l := range links {
			nl, err := url.NewLink(l.Href, &lnk)
			if err != nil {
				fmt.Printf("%s link parse error: %s\n", lnk, err)
				fmt.Println(err)
				continue
			}
			nl.Unify(lnk)

			if nl.Host != lnk.Host {
				continue
			}

			j.OutChan <- nl
		}
		<-j.SyncChan
	}
}
