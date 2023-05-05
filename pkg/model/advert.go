package model

type SortBy int
type SortOrder int

const (
	Price SortBy = iota
	Date
)

const (
	Asc SortOrder = iota
	Desc
)

type DetailedAdvert struct {
	Id          int64    `json:"id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Photos      []string `json:"photos"`
	Price       float64  `json:"price"`
}

type Advert struct {
	Id        int64   `json:"id"`
	Title     string  `json:"title"`
	MainPhoto string  `json:"main_photo"`
	Price     float64 `json:"price"`
}
