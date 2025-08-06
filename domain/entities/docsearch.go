package entities

type SerpAPIResponse struct {
	SearchParameters struct {
		Q      string `json:"q"`
		Engine string `json:"engine"`
	} `json:"search_parameters"`
	OrganicResults []struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
	} `json:"organic_results"`
	// Add other fields you might need, e.g., "knowledge_graph", "related_searches", etc.
}
