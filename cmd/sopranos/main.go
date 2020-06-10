package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/docker/docker/client"

	"github.com/liampulles/rikyu/pkg/dvd"
	"github.com/liampulles/rikyu/pkg/exec"
	"github.com/liampulles/rikyu/pkg/system"
)

// TODO: What about a ripping DSL? Think about how that would factor into the stages of i.e. ripping,
// creating filters, creating pipelines (chained builder?), immutable (run again and resume),
// - A script for each dvd rip
// - A script for every episode encode
// - A script for general category encodes, e.g. deleted scenes.

func main() {
	ds := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return client.NewEnvClient()
	})
	ss := system.NewSystemServiceImpl(
		os.MkdirAll, os.Stat,
		ioutil.ReadDir, ioutil.WriteFile, ioutil.ReadFile)
	dvdS := dvd.NewDVDServiceImpl(ds, ss)

	out := "/media/liam/Additional 2/Encodes/Sopranos1/"
	err := dvdS.RipTitles("/media/liam/SOPRANOS_SEASON1_DISC1/", map[int]string{
		1: path.Join(out, "episodes/1"),
		3: path.Join(out, "episodes/2"),
		4: path.Join(out, "episodes/3"),
	})
	if err != nil {
		panic(err)
	}
}
