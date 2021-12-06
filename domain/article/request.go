package article

// CreateArticleRequest is model for creating article.
type CreateArticleRequest struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

// EditArticleRequest is model for modified article.
type EditArticleRequest struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}
