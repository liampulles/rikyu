package dvd

import (
	"encoding/xml"
	"math"
	"time"

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
	ReadDVDInfoForDirectory(dir string) (*DVD, error)
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

func (dis *DVDInfoServiceImpl) ReadDVDInfoForDirectory(dir string) (*DVD, error) {
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

func mapLsdvdStructToDVD(in *lsdvdDVD) *DVD {
	titleSetMap := make(map[int][]Title)
	for _, track := range in.Track {
		// Audio
		var audioTracks []Audio
		for _, lsAud := range track.Audio {
			audioTrack := Audio{
				Language: Language{
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
		var subtitleTracks []Subtitle
		for _, lsSub := range track.Subp {
			subtitleTrack := Subtitle{
				Language: Language{
					Code: lsSub.LangCode,
					Name: lsSub.Language,
				},
				Content: lsSub.Content,
			}
			subtitleTracks = append(subtitleTracks, subtitleTrack)
		}
		// Chapters
		var chapters []Chapter
		for _, lsChap := range track.Chapter {
			chapter := Chapter{
				Length:    durationFromSecs(lsChap.Length),
				StartCell: lsChap.StartCell,
			}
			chapters = append(chapters, chapter)
		}
		// Cells
		length := time.Duration(0)
		var cells []Cell
		for _, lsCell := range track.Cell {
			cell := Cell{
				Length: durationFromSecs(lsCell.Length),
			}
			cells = append(cells, cell)
			length += cell.Length
		}
		// Palette
		var palette Palette
		for i, lsCol := range track.Palette {
			palette[i] = lsCol
		}
		// Finally, the title
		title := Title{
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
	var titleSets []TitleSet
	for _, titles := range titleSetMap {
		titleSets = append(titleSets, TitleSet{
			Titles: titles,
		})
	}
	return &DVD{
		DiscTitle: in.Title,
		TitleSets: titleSets,
	}
}

func durationFromSecs(secs float64) time.Duration {
	actualSecs := math.Floor(secs)
	milliSecs := (secs - actualSecs) * 1000
	return time.Duration(float64(time.Second)*actualSecs) + time.Duration(float64(time.Millisecond)*milliSecs)
}
