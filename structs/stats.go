package structs

type Stats struct {
	PostViews         map[int]int
	SearchAppearances map[int]int
	AuthorViews       map[string]int
	TotalSearchCount  int
}

var GlobalStats = Stats{
	PostViews:         make(map[int]int),
	SearchAppearances: make(map[int]int),
	AuthorViews:       make(map[string]int),
	TotalSearchCount:  0,
}
