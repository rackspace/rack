package objectcommands

import (
	"flag"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/output"
)

func newUpDirCmd(fs *flag.FlagSet) *commandUploadDir {
	return &commandUploadDir{Ctx: &handler.Context{
		CLIContext: cli.NewContext(cli.NewApp(), fs, nil),
	}}
}

func TestUploadDirContext(t *testing.T) {
	cmd := newUpDirCmd(flag.NewFlagSet("flags", 1))
	th.AssertDeepEquals(t, cmd.Ctx, cmd.Context())
}

func TestUploadDirKeys(t *testing.T) {
	cmd := &commandUploadDir{}
	th.AssertDeepEquals(t, keysUploadDir, cmd.Keys())
}

func TestUploadDirServiceClientType(t *testing.T) {
	cmd := &commandUploadDir{}
	th.AssertEquals(t, serviceClientType, cmd.ServiceClientType())
}

func TestUploadDirErrWhenCtnrMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)

	err := newUpDirCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--container is required."}
	th.AssertDeepEquals(t, expected, err)
}

/*
func TestWarningEmittedForNonDirs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	fs := flag.NewFlagSet("flags", 1)

	_, filename, _, _ := runtime.Caller(0)

	fs.String("container", "", "")
	fs.String("dir", "", "")
	fs.Set("container", "foo")
	fs.Set("dir", filename)

	cmd := newUpDirCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{}
	cmd.HandleFlags(res)
	cmd.Execute(res)

	fmt.Println(filename)

	err := fmt.Errorf("%s is not a directory, ignoring", filename)
	th.AssertDeepEquals(t, err, res.Err)
}
*/

/*
func TestRecursionIsDisabledByDefault(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	_, filename, _, _ := runtime.Caller(0)
	rootDirFix := path.Dir(path.Dir(filename))

	th.Mux.HandleFunc("/foo/commands.go", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		str, err := ioutil.ReadFile(rootDirFix + "/commands.go")
		if err == nil {
			th.TestBody(t, r, string(str))
		}
		w.WriteHeader(201)
	})

	fs := flag.NewFlagSet("flags", 1)

	fs.String("container", "", "")
	fs.String("dir", "", "")
	fs.String("quiet", "", "")

	fs.Set("container", "foo")
	fs.Set("dir", rootDirFix)
	fs.Set("quiet", "true")

	cmd := newUpDirCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{}
	cmd.HandleFlags(res)
	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)

	if !strings.Contains(res.Result.(string), "Uploaded 1 object ") {
		t.Fatalf("Unexpected result message: %s", res.Result)
	}
}

func TestUploadDirExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	_, filename, _, _ := runtime.Caller(0)
	rootDirFix := path.Dir(filename)

	var count int64
	count = 0

	filepath.Walk(rootDirFix, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		urlPath := "/foo" + strings.TrimPrefix(path, rootDirFix)
		th.Mux.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			str, err := ioutil.ReadFile(path)
			if err == nil {
				th.TestBody(t, r, string(str))
				count++
			}
			w.WriteHeader(201)
		})
		return nil
	})

	fs := flag.NewFlagSet("flags", 1)

	fs.String("container", "", "")
	fs.String("dir", "", "")
	fs.String("quiet", "", "")

	fs.Set("container", "foo")
	fs.Set("dir", rootDirFix)
	fs.Set("quiet", "true")

	cmd := newUpDirCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{}
	cmd.HandleFlags(res)
	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)

	if !strings.Contains(res.Result.(string), fmt.Sprintf("Uploaded %d %s ", count, util.Pluralize("object", count))) {
		t.Fatalf("Unexpected result message: %s", res.Result)
	}
}
*/
