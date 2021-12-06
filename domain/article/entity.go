package article

import "time"

// ArticleStatus is a type of article current status.
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "DRAFT"
	ArticleStatusPublished ArticleStatus = "PUBLISHED"
)

// Article is a collection of property of article.
type Article struct {
	ID             int64
	Title          string
	Subtitle       string
	Content        string
	Status         ArticleStatus
	CreatedAt      time.Time
	PublishedAt    *time.Time
	LastModifiedAt *time.Time
}
