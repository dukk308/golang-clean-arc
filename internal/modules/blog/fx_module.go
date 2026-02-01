package blog

import (
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/application"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/infrastructure/persistence"
	blog_http "github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/presentation/http"
	"go.uber.org/fx"
)

var Module = fx.Module("blog",
	fx.Provide(
		fx.Annotate(
			persistence.NewBlogRepository,
			fx.As(new(domain.IBlogRepository)),
		),
	),
	fx.Provide(application.NewCreateBlogCommand),
	fx.Provide(application.NewGetBlogQuery),
	fx.Provide(application.NewListBlogsQuery),
	fx.Provide(application.NewUpdateBlogCommand),
	fx.Provide(application.NewDeleteBlogCommand),
	fx.Provide(
		func(
			createBlogCommand *application.CreateBlogCommand,
			getBlogQuery *application.GetBlogQuery,
			listBlogsQuery *application.ListBlogsQuery,
			updateBlogCommand *application.UpdateBlogCommand,
			deleteBlogCommand *application.DeleteBlogCommand,
		) *blog_http.Http {
			return blog_http.NewHttp(
				createBlogCommand,
				getBlogQuery,
				listBlogsQuery,
				updateBlogCommand,
				deleteBlogCommand,
			)
		},
	),
)
