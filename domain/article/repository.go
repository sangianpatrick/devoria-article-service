package article

import "context"

type ArticleRepository interface {
	Save(ctx context.Context, article Article) (ID int64, err error)
	Update(ctx context.Context, ID int64, updatedArticle Article) (err error)
	FindByID(ctx context.Context, ID int64) (article Article, err error)
	FindMany(ctx context.Context) (bunchOfArticles []Article, err error)
	FindManySpecificProfile(ctx context.Context, articleID int64) (bunchOfArticles []Article, err error)
}
