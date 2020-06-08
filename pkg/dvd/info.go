package dvd

import (
	"encoding/xml"
	"math"
	"time"

	"github.com/liampulles/rikyu/pkg/exec"
	"github.com/liampulles/rikyu/pkg/types"
)

func (dis *DVDServiceImpl) ReadDVDInfoForDirectory(dir string) (*types.DVD, error) {
	// Get StdOut
	mounts := []exec.DockerVolumeMount{
		{
			Host:      dir,
			Container: "/dvd",
		},
	}
	args := []string{
		"lsdvd",
		"-qx",
		"-Ox",
		"/dvd",
	}
	stdOut, _, err := dis.dockerService.RunDockerContainerForOutput(lsdvdImage, mounts, args)
	if err != nil {
		return nil, err
	}

	// Unmarshal into lsdvd structs
	var data lsdvdDVD
	if err := xml.Unmarshal([]byte(stdOut), &data); err != nil {
		return nil, err
	}

	return mapLsdvdStructToDVD(&data), nil
}

func mapLsdvdStructToDVD(in *lsdvdDVD) *types.DVD {
	titleSetMap := make(map[int][]types.Title)
	for _, track := range in.Track {
		// Audio
		var audioTracks []types.Audio
		for _, lsAud := range track.Audio {
			audioTrack := types.Audio{
				Language: types.Language{
					Code: lsAud.LangCode,
					Name: lsAud.Language,
				},
				Content:             lsAud.Content,
				Format:              lsAud.Format,
				Frequency:           lsAud.Frequency,
				DynamicRangeControl: lsAud.Quantization == "drc",
				Channels:            lsAud.Channels,
			}
			audioTracks = append(audioTracks, audioTrack)
		}
		// Subtitles
		var subtitleTracks []types.Subtitle
		for _, lsSub := range track.Subp {
			subtitleTrack := types.Subtitle{
				Language: types.Language{
					Code: lsSub.LangCode,
					Name: lsSub.Language,
				},
				Content: lsSub.Content,
			}
			subtitleTracks = append(subtitleTracks, subtitleTrack)
		}
		// Chapters
		var chapters []types.Chapter
		for _, lsChap := range track.Chapter {
			chapter := types.Chapter{
				Length:    durationFromSecs(lsChap.Length),
				StartCell: lsChap.StartCell,
			}
			chapters = append(chapters, chapter)
		}
		// Cells
		length := time.Duration(0)
		var cells []types.Cell
		for _, lsCell := range track.Cell {
			cell := types.Cell{
				Length: durationFromSecs(lsCell.Length),
			}
			cells = append(cells, cell)
			length += cell.Length
		}
		// Palette
		var palette types.Palette
		for i, lsCol := range track.Palette {
			palette[i] = lsCol
		}
		// Finally, the title
		title := types.Title{
			TitleNumber:        track.IX,
			Angles:             track.Angles,
			FPS:                track.FPS,
			Width:              track.Width,
			Height:             track.Height,
			Format:             track.Format,
			AutomaticLetterbox: track.DF == "Letterbox",
			Length:             length,
			AudioTracks:        audioTracks,
			SubtitleTracks:     subtitleTracks,
			Chapters:           chapters,
			Cells:              cells,
			Palette:            palette,
		}
		titleSetMap[track.Vts] = append(titleSetMap[track.Vts], title)
	}
	// TitleSets
	var titleSets []types.TitleSet
	for _, titles := range titleSetMap {
		titleSets = append(titleSets, types.TitleSet{
			Titles: titles,
		})
	}
	return &types.DVD{
		DiscTitle: in.Title,
		TitleSets: titleSets,
	}
}

func durationFromSecs(secs float64) time.Duration {
	actualSecs := math.Floor(secs)
	milliSecs := (secs - actualSecs) * 1000
	return time.Duration(float64(time.Second)*actualSecs) + time.Duration(float64(time.Millisecond)*milliSecs)
}
