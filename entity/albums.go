package entity

type Albums struct {
	Sets Photosets `json:"photosets"`
	Stat string    `json:"stat"`
}

type Photosets struct {
	Pages    int64   `json:"pages"`
	Total    int64   `json:"total"`
	Photoset []Album `json:"photoset"`
}

type Album struct {
	ID          string  `json:"id"`
	Title       Content `json:"title"`
	Description Content `json:"description"`
}

type Content struct {
	Content string `json:"_content"`
}
