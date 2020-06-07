package dvd

import (
	"time"
)

type ChannelFormat string

type FPS struct {
	Numerator   int
	Denominator int
	Value       float64
}

type Language struct {
	Code string
	Name string
}

type AspectRation struct {
	Numerator   int
	Denominator int
	Value       float64
}

type DVD struct {
	DiscTitle string
	TitleSets []TitleSet
}

type TitleSet struct {
	Titles []Title
}

type Title struct {
	TitleNumber        int
	Angles             int
	FPS                float64
	Width              int
	Height             int
	Format             string
	AutomaticLetterbox bool
	Length             time.Duration
	AudioTracks        []Audio
	SubtitleTracks     []Subtitle
	Chapters           []Chapter
	Cells              []Cell
	Palette            Palette
}

type Audio struct {
	Language            Language
	Content             string
	Format              string
	Frequency           int
	DynamicRangeControl bool
	Channels            int
}

type Subtitle struct {
	Language Language
	Content  string
}

type Chapter struct {
	Length    time.Duration
	StartCell int
}

type Cell struct {
	Length time.Duration
}

type Palette [16]string
