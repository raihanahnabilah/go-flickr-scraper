package entity

type Links struct {
	Sizes Sizes  `json:"sizes"`
	Stat  string `json:"stat"`
}

type Sizes struct {
	Size []Size `json:"size"`
}

type Size struct {
	Label  string `json:"label"`
	Source string `json:"source"`
}
