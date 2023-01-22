package main

import (
	"flag"
	"fmt"
	"os"

	job "github.com/volli1704/sitemap/internal/job"
	"github.com/volli1704/sitemap/internal/linkmap"
	"github.com/volli1704/sitemap/internal/url"
)

var lm linkmap.LinkMap = make(linkmap.LinkMap)

func main() {
	lnk := flag.String("url", "https://github.com/", "Root url of the site from where app will explore link tree")
	depth := flag.Int("depth", 2, "Depth of link research")
	workerCount := flag.Int("workers", 50, "Count of workers to parse site urls")
	outFile := flag.String("o", "result.out", "Output file")
	flag.Parse()

	f, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	root, err := url.NewLink(*lnk, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("start goroutine")
	jr := job.NewJobRunner(int8(*workerCount), *depth)
	jr.RunWorkers()
	jr.Send(root)
	jr.Start(lm)

	f.Write([]byte(lm.String()))
}
