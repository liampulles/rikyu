package types

type Project struct {
	Name string `json:"name"`
	// TODO: Templates
	// TODO: Jobs
	Path string `json:"-"`
}

type UniqueCategory struct {
	Name  string `json:"name"`
	Title Title  `json:"title"`
}

type SequentialCategory struct {
	Name   string        `json:"name"`
	Titles map[int]Title `json:"titles"`
}
