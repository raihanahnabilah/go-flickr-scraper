package entity

type Photos struct {
	Photoset Photoset `json:"photoset"`
	Stat     string   `json:"stat"`
}

type Photoset struct {
	ID    string  `json:"id"`
	Photo []Photo `json:"photo"`
	Pages int64   `json:"pages"`
}

type Photo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
