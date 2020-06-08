package dvd

import (
	"encoding/xml"
	"math"
	"time"

	"github.com/liampulles/rikyu/pkg/types"

	"github.com/liampulles/rikyu/pkg/exec"
)

const (
	lsdvdImage = "docker.io/lpulles/dvdtools:latest"
)

type lsdvdDVD struct {
	Device     string       `xml:"device"`
	Title      string       `xml:"title"`
	VmgID      string       `xml:"vmg_id"`
	ProviderID string       `xml:"provider_id"`
	Track      []lsdvdTrack `xml:"track,omitempty"`
}

type lsdvdTrack struct {
	IX      int            `xml:"ix"`
	Length  float64        `xml:"length"`
	VtsID   string         `xml:"vts_id"`
	Vts     int            `xml:"vts"`
	Ttn     int            `xml:"ttn"`
	FPS     float64        `xml:"fps"`
	Format  string         `xml:"format"`
	Aspect  string         `xml:"aspect"`
	Width   int            `xml:"width"`
	Height  int            `xml:"height"`
	DF      string         `xml:"df"`
	Palette []string       `xml:"palette>color,omitempty"`
	Angles  int            `xml:"angles"`
	Audio   []lsdvdAudio   `xml:"audio,omitempty"`
	Chapter []lsdvdChapter `xml:"chapter,omitempty"`
	Cell    []lsdvdCell    `xml:"cell,omitempty"`
	Subp    []lsdvdSubp    `xml:"subp,omitempty"`
}

type lsdvdAudio struct {
	IX           int    `xml:"ix"`
	LangCode     string `xml:"langcode"`
	Language     string `xml:"language"`
	Format       string `xml:"format"`
	Frequency    int    `xml:"frequency"`
	Quantization string `xml:"quantization"`
	Channels     int    `xml:"channels"`
	ApMode       int    `xml:"ap_mode"`
	Content      string `xml:"content"`
	StreamID     string `xml:"streamid"`
}

type lsdvdChapter struct {
	IX        int     `xml:"ix"`
	Length    float64 `xml:"length"`
	StartCell int     `xml:"startcell"`
}

type lsdvdCell struct {
	IX     int     `xml:"ix"`
	Length float64 `xml:"length"`
}

type lsdvdSubp struct {
	IX       int    `xml:"ix"`
	LangCode string `xml:"langcode"`
	Language string `xml:"language"`
	Content  string `xml:"content"`
	StreamID string `xml:"streamid"`
}

type DVDInfoService interface {
	ReadDVDInfoForDirectory(dir string) (*types.DVD, error)
}

type DVDInfoServiceImpl struct {
	dockerService exec.DockerService
}

var _ DVDInfoService = &DVDInfoServiceImpl{}

func NewDVDInfoServiceImpl(dockerService exec.DockerService) *DVDInfoServiceImpl {
	return &DVDInfoServiceImpl{
		dockerService: dockerService,
	}
}

func (dis *DVDInfoServiceImpl) ReadDVDInfoForDirectory(dir string) (*types.DVD, error) {
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
