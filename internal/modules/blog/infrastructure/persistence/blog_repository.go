package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) domain.IBlogRepository {
	return &BlogRepository{
		db: db,
	}
}

func (r *BlogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	sqlBlog := &SQLBlog{}
	sqlBlog.FromDomain(blog)
	return r.db.WithContext(ctx).Create(sqlBlog).Error
}

func (r *BlogRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&SQLBlog{}).Error
}

func (r *BlogRepository) GetPage(ctx context.Context, offset, limit int) ([]*domain.Blog, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&SQLBlog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var sqlBlogs []SQLBlog
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&sqlBlogs).Error; err != nil {
		return nil, 0, err
	}
	blogs := make([]*domain.Blog, len(sqlBlogs))
	for i, b := range sqlBlogs {
		blogs[i] = b.ToDomain()
	}
	return blogs, total, nil
}

func (r *BlogRepository) GetBySlug(ctx context.Context, slug string) (*domain.Blog, error) {
	var sqlBlog SQLBlog
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&sqlBlog).Error; err != nil {
		return nil, err
	}
	return sqlBlog.ToDomain(), nil
}

func (r *BlogRepository) GetByID(ctx context.Context, id string) (*domain.Blog, error) {
	var sqlBlog SQLBlog
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&sqlBlog).Error; err != nil {
		return nil, err
	}
	return sqlBlog.ToDomain(), nil
}

func (r *BlogRepository) Update(ctx context.Context, blog *domain.Blog) error {
	sqlBlog := &SQLBlog{}
	sqlBlog.FromDomain(blog)
	return r.db.WithContext(ctx).Save(sqlBlog).Error
}
