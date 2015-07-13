package objectcommands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/cenkalti/backoff"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/dustin/go-humanize"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/jrperritt/rack/util"
)

var uploadDir = cli.Command{
	Name:        "upload-dir",
	Usage:       util.Usage(commandPrefix, "upload-dir", "--container <containerName> --dir <dirName>"),
	Description: "Uploads the contents of a local directory to a container",
	Action:      actionUploadDir,
	Flags:       util.CommandFlags(flagsUploadDir, keysUploadDir),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUploadDir, keysUploadDir))
	},
}

func flagsUploadDir() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container to upload the object upload",
		},
		cli.StringFlag{
			Name:  "dir",
			Usage: "[required] The name the local directory which will be uploaded",
		},
		cli.StringFlag{
			Name:  "content-type",
			Usage: "[optional] The Content-Type header that will be set on all objects.",
		},
		cli.IntFlag{
			Name:  "content-length",
			Usage: "[optional] The Content-Length header that will be set on all objects.",
		},
		cli.StringFlag{
			Name:  "content-encoding",
			Usage: "[optional] The Content-Encoding header that will be set on all objects. By default, the uploaded content will be gzipped.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `dir` isn't provided] The field being piped to STDIN. Valid values are: dir",
		},
		cli.IntFlag{
			Name:  "concurrency",
			Usage: "[optional] The amount of concurrent workers that will upload the directory.",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "[optional] By default every file upload will be outputted. If --quiet is provided, only a final summary will be outputted.",
		},
	}
}

var keysUploadDir = []string{}

type paramsUploadDir struct {
	container   string
	dir         string
	stream      io.ReadSeeker
	opts        objects.CreateOpts
	concurrency int
	quiet       bool
}

type commandUploadDir handler.Command

func actionUploadDir(c *cli.Context) {
	command := &commandUploadDir{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUploadDir) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUploadDir) Keys() []string {
	return keysUploadDir
}

func (command *commandUploadDir) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUploadDir) HandleFlags(resource *handler.Resource) error {
	if err := command.Ctx.CheckFlagsSet([]string{"container", "dir"}); err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := objects.CreateOpts{
		ContentLength: int64(c.Int("content-length")),
		ContentType:   c.String("content-type"),
	}

	if c.IsSet("content-encoding") && c.String("content-encoding") != "gzip" {
		opts.ContentEncoding = c.String("content-encoding")
	}

	conc := c.Int("concurrency")
	if conc <= 0 {
		conc = 5
	}

	resource.Params = &paramsUploadDir{
		container:   c.String("container"),
		dir:         c.String("dir"),
		opts:        opts,
		concurrency: c.Int("concurrency"),
		quiet:       c.Bool("quiet"),
	}

	return nil
}

func (command *commandUploadDir) StdinField() string {
	return "dir"
}

func (command *commandUploadDir) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsUploadDir).dir = item
	return nil
}

func (command *commandUploadDir) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"dir"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsUploadDir).dir = command.Ctx.CLIContext.String("dir")
	return nil
}

func (command *commandUploadDir) Execute(resource *handler.Resource) {
	// bump thread count to number of available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	params := resource.Params.(*paramsUploadDir)
	jobs := make(chan string)
	results := make(chan *handler.Resource)

	var wg sync.WaitGroup
	var totalSize uint64
	start := time.Now()

	go func() {
		for r := range results {
			command.Ctx.WaitGroup.Add(1)
			command.Ctx.Results <- r
		}
		command.Ctx.WaitGroup.Wait()
	}()

	for i := 0; i < params.concurrency; i++ {
		wg.Add(1)
		go func(totalSize *uint64) {
			for p := range jobs {
				var re *handler.Resource

				ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
				for _ = range ticker.C {
					re = command.handle(p, params)
					if re.Err != nil {
						continue
					}

					ticker.Stop()
					break
				}

				fi, err := os.Stat(p)
				if err == nil {
					*totalSize += uint64(fi.Size())
				}

				results <- re
			}
			wg.Done()
		}(&totalSize)
	}

	filepath.Walk(params.dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			jobs <- path
		}
		return nil
	})
	close(jobs)

	wg.Wait()
	close(results)

	resource.Result = fmt.Sprintf("Finished! %s uploaded in %s\n", humanize.Bytes(totalSize), humanize.RelTime(start, time.Now(), "", ""))
}

func (command *commandUploadDir) handle(p string, params *paramsUploadDir) *handler.Resource {
	re := &handler.Resource{}

	file, err := os.Open(p)
	defer file.Close()

	if err != nil {
		re.Err = err
		return re
	}

	on := strings.TrimPrefix(p, params.dir+"/")
	res := objects.Create(command.Ctx.ServiceClient, params.container, on, file, params.opts)
	re.Err = res.Err

	if res.Err == nil {
		if params.quiet == true {
			re.Result = ""
		} else {
			re.Result = fmt.Sprintf("Uploaded %s to %s\n", on, params.container)
		}
	}

	return re
}
