package job

import (
	"time"

	"github.com/volli1704/sitemap/internal/linkmap"
	url "github.com/volli1704/sitemap/internal/url"
)

// Runs, stops and synchronize url parser workers
type runner struct {
	InChan      chan url.URL
	OutChan     chan url.URL
	SyncChan    chan bool
	WorkerCount int8
	Workers     []worker
	MaxDepth    int
}

// Every x seconds job runner will check if there're any workers running, if not it will stop `url loop`
const tickDuration = 1 * time.Second

// Initialize JobRunner with worker count and depth of page research
func NewJobRunner(workerCount int8, maxDepth int) runner {
	in := make(chan url.URL, 1)
	out := make(chan url.URL, 1<<16)
	sync := make(chan bool, 100)

	return runner{
		in,
		out,
		sync,
		workerCount,
		make([]worker, 0),
		maxDepth,
	}
}

// Start `url loop`. When worker parse page, all urls on this page go to out channel,
// runner checks if these url were processed and if not, push it to the inner channel and
// workers will process it in a `loop`
func (r *runner) Start(res linkmap.LinkMap) {
	ticker := time.NewTicker(tickDuration)
lp:
	for {
		select {
		case out := <-r.OutChan:
			if out.Depth > r.MaxDepth {
				continue
			}

			if _, ok := res[out.String()]; !ok {
				res[out.String()] = 0
				r.InChan <- out
			}
			res[out.String()]++
		case <-ticker.C:
			if len(r.SyncChan) == 0 {
				//fmt.Println(res.String())

				ticker.Stop()
				break lp
			}
		}
	}
}

// Run all workers
func (r *runner) RunWorkers() {
	for i := int8(1); i <= r.WorkerCount; i++ {
		w := NewWorker(i, r.InChan, r.OutChan, r.SyncChan)
		r.Workers = append(r.Workers, w)

		go w.run()
	}
}

// Send url to inner loop. In most cases we need it to send initial URL to
// the inner channel
func (r *runner) Send(url url.URL) {
	r.InChan <- url
}
