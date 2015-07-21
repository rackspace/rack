package containercommands

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/cenkalti/backoff"
	"github.com/jrperritt/rack/internal/github.com/dustin/go-humanize"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/pagination"
	"github.com/jrperritt/rack/util"
)

type handleEmptyParams struct {
	container   string
	quiet       bool
	concurrency int
}

func handleEmpty(command handler.Commander, resource *handler.Resource, params *handleEmptyParams) {
	// bump thread count to number of available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	jobs := make(chan string)

	var wg sync.WaitGroup
	var totalFiles int64
	start := time.Now()

	if params.concurrency == 0 {
		params.concurrency = 100
	}

	for i := 0; i < params.concurrency; i++ {
		wg.Add(1)
		go func(totalFiles *int64) {
			for objectName := range jobs {
				re := &handler.Resource{}

				ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
				for _ = range ticker.C {
					rawResponse := objects.Delete(command.Context().ServiceClient, params.container, objectName, nil)
					if rawResponse.Err != nil {
						continue
					}
					re.Result = fmt.Sprintf("Successfully deleted object [%s] from container [%s]\n", objectName, params.container)

					ticker.Stop()
					break
				}

				*totalFiles++

				if !params.quiet {
					command.Context().Results <- re
				}
			}
			wg.Done()
		}(&totalFiles)
	}

	pager := objects.List(command.Context().ServiceClient, params.container, nil)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		names, err := objects.ExtractNames(page)
		if err != nil {
			return false, err
		}
		for _, name := range names {
			jobs <- name
		}
		return true, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
	close(jobs)

	wg.Wait()
	resource.Result = fmt.Sprintf("Finished! Deleted %s %s in %s", humanize.Comma(totalFiles), util.Pluralize("object", totalFiles), humanize.RelTime(start, time.Now(), "", ""))
}
