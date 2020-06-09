package dvd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liampulles/rikyu/pkg/exec"
)

func (dsi *DVDServiceImpl) RipTitles(dvdDir string, titlesToDirs map[int]string) error {
	for titleNo, outDir := range titlesToDirs {
		titleStr := strconv.Itoa(titleNo)
		if err := dsi.systemService.CreateDirRecursive(outDir); err != nil {
			return err
		}
		if err := dsi.validateDirIsEmpty(outDir); err != nil {
			return err
		}
		mounts := []exec.DockerVolumeMount{
			{
				Host:      dvdDir,
				Container: "/in",
			},
			{
				Host:      outDir,
				Container: "/out",
			},
		}
		args := []string{
			"/sh/extractTitle.sh",
			titleStr,
		}
		if _, _, err := dsi.dockerService.RunDockerContainerForOutput(dvdbackupImage, mounts, args); err != nil {
			return err
		}
		if err := dsi.validateDirHasVobAndIfo(outDir); err != nil {
			return err
		}
	}
	return nil
}

func (dsi *DVDServiceImpl) validateDirIsEmpty(dir string) error {
	files, err := dsi.systemService.GetFilesInDir(dir)
	if err != nil {
		return err
	}
	if len(files) != 0 {
		return fmt.Errorf("expected an empty dir but say some files: %v", files)
	}
	return nil
}

func (dsi *DVDServiceImpl) validateDirHasVobAndIfo(dir string) error {
	files, err := dsi.systemService.GetFilesInDir(dir)
	if err != nil {
		return err
	}

	var containsVob, containsIfo bool
	for _, file := range files {
		upper := strings.ToUpper(file)
		if strings.HasSuffix(upper, "VOB") {
			containsVob = true
		} else if strings.HasSuffix(upper, "IFO") {
			containsIfo = true
		}
		if containsVob && containsIfo {
			break
		}
	}

	if !containsVob {
		return fmt.Errorf("expected dir to contain VOB, but no such file was found")
	}
	if !containsIfo {
		return fmt.Errorf("expected dir to contain IFO, but no such file was found")
	}

	return nil
}
