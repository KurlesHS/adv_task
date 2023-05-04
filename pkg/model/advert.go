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
	Id          int64
	Title       string
	Description string
	Photos      []string
	Price       float64
}

type Advert struct {
	Id        int64
	Title     string
	MainPhoto string
	Price     float64
}
