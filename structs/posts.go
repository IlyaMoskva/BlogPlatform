package structs

type ReportResponse struct {
	TotalViews        int      `json:"total_views"`
	SearchAppearances int      `json:"search_appearances"`
	TopAuthors        []Author `json:"top_authors"`
	TopPosts          []Post   `json:"top_posts"`
}

type Author struct {
	Name  string `json:"name"`
	Views int    `json:"views"`
}

type Post struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	Author            string `json:"author"`
	Views             int    `json:"views"`
	SearchAppearances int    `json:"search_appearances"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
