package entities

type SerpAPIResponse struct {
	SearchParameters SearchParameters `json:"search_parameters"`
	OrganicResults   OrganicResult    `json:"organic_results"`
	// Add other fields you might need, e.g., "knowledge_graph", "related_searches", etc.
}

type SerpURL struct {
	URL         string  `json:"url"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

type SerpURLScore struct {
	SerpURL []SerpURL `json:"serpURL"`
}

type SerpReturn struct {
	RefId   string        `json:"refId"`
	Content string        `json:"content"`
	Source  OrganicResult `json:"source"`
}
type OrganicResult []struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

type SearchParameters struct {
	Q      string `json:"q"`
	Engine string `json:"engine"`
}
