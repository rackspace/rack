package containercommands

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/rackspace/rack/handler"
	"github.com/cenkalti/backoff"
	"github.com/dustin/go-humanize"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/rack/util"
)

type handleEmptyParams struct {
	container   string
	quiet       bool
	concurrency int
}

func handleEmpty(command handler.Commander, resource *handler.Resource, params *handleEmptyParams) {
	var totalFiles int64
	var wg sync.WaitGroup

	// bump thread count to number of available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// get the names of all the objects in the container
	allPages, err := objects.List(command.Context().ServiceClient, params.container, nil).AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	names, err := objects.ExtractNames(allPages)
	if err != nil {
		resource.Err = err
		return
	}

	// send the object names into the `jobs` channel
	jobs := make(chan string, len(names))
	for i := 0; i < len(names); i++ {
		wg.Add(1)
		jobs <- names[i]
	}
	close(jobs)

	// default the number of goroutines to spawn if the `concurrency` flag
	// wasn't provided
	if params.concurrency == 0 {
		params.concurrency = 100
	}

	start := time.Now()

	for i := 0; i < params.concurrency; i++ {
		go func(totalFiles *int64) {
			for objectName := range jobs {
				ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
				for _ = range ticker.C {
					rawResponse := objects.Delete(command.Context().ServiceClient, params.container, objectName, nil)
					if rawResponse.Err != nil {
						continue
					}
					ticker.Stop()
					break
				}

				*totalFiles++

				if !params.quiet {
					re := &handler.Resource{
						Result: fmt.Sprintf("Successfully deleted object [%s] from container [%s]\n", objectName, params.container),
					}
					command.Context().Results <- re
				}
				wg.Done()
			}
		}(&totalFiles)
	}

	wg.Wait()

	resource.Result = fmt.Sprintf("Finished! Deleted %s %s in %s", humanize.Comma(totalFiles), util.Pluralize("object", totalFiles), humanize.RelTime(start, time.Now(), "", ""))
}
