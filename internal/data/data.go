package data

type Book struct {
	ID       int64   `json:"id"`
	Title    string  `json:"title"`
	Language *string `json:"language,omitempty"`
	PubYear  *int    `json:"pub_year,omitempty"`
	ISBN     *string `json:"isbn,omitempty"`
}

type BookWithMeta struct {
	Book
	Authors   []string `json:"authors"`
	Genres    []string `json:"genres"`
	AvgRating *float64 `json:"avg_rating,omitempty"`
}
