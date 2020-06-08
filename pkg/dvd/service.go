package dvd

import (
	"github.com/liampulles/rikyu/pkg/exec"
	"github.com/liampulles/rikyu/pkg/types"
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

type DVDService interface {
	ReadDVDInfoForDirectory(dir string) (*types.DVD, error)
}

type DVDServiceImpl struct {
	dockerService exec.DockerService
}

var _ DVDService = &DVDServiceImpl{}

func NewDVDServiceImpl(dockerService exec.DockerService) *DVDServiceImpl {
	return &DVDServiceImpl{
		dockerService: dockerService,
	}
}
