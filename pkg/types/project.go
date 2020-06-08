package types

type Project struct {
	UniqueCategories     []UniqueCategory
	SequentialCategories []SequentialCategory
	// TODO: Templates
	// TODO: Jobs
	// TODO: Path on hard drive
}

type UniqueCategory struct {
	Name  string
	Title Title
}

type SequentialCategory struct {
	Name   string
	Titles map[int]Title
}
