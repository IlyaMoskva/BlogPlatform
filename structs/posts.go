package structs

type Post struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	Author            string `json:"author"`
	Views             int    `json:"views"`
	SearchAppearances int    `json:"search_appearances"`
}
