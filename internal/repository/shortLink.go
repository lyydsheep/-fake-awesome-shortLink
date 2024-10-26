package repository

import (
	"awesome-shortLink/internal/domain"
	"awesome-shortLink/internal/repository/cache"
	"awesome-shortLink/internal/repository/dao"
	"context"
	"go.uber.org/zap"
)

var (
	ErrDuplicatedLongURL = dao.ErrDuplicatedLongURL
	ErrNotFound          = dao.ErrNotFound
)

type ShortLinkRepository interface {
	Create(ctx context.Context, longURL string) (domain.ShortLink, error)
	FindByShort(ctx context.Context, shortURL string) (domain.ShortLink, error)
}

type ShortLinkRepositoryV1 struct {
	dao dao.ShortLinkDAO
	l   *zap.Logger
}

func (repo *ShortLinkRepositoryV1) FindByShort(ctx context.Context, shortURL string) (domain.ShortLink, error) {
	sl, err := repo.dao.FindByShort(ctx, shortURL)
	if err != nil {
		return domain.ShortLink{}, err
	}
	return entityToDomain(sl), nil
}

func (repo *ShortLinkRepositoryV1) Create(ctx context.Context, url string) (domain.ShortLink, error) {
	sl, err := repo.dao.Insert(ctx, url)
	if err != nil {
		return domain.ShortLink{}, err
	}
	return entityToDomain(sl), nil
}

func NewShortLinkRepositoryV1(dao dao.ShortLinkDAO, l *zap.Logger) ShortLinkRepository {
	return &ShortLinkRepositoryV1{
		dao: dao,
		l:   l,
	}
}

type ShortLinkRepositoryV2 struct {
	dao   dao.ShortLinkDAO
	l     *zap.Logger
	cache cache.ShortLinkCache
}

func (repo *ShortLinkRepositoryV2) Create(ctx context.Context, longURL string) (domain.ShortLink, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *ShortLinkRepositoryV2) FindByShort(ctx context.Context, shortURL string) (domain.ShortLink, error) {

}

func NewShortLinkRepositoryV2(dao dao.ShortLinkDAO, l *zap.Logger, cache cache.ShortLinkCache) ShortLinkRepository {
	return &ShortLinkRepositoryV2{
		dao:   dao,
		l:     l,
		cache: cache,
	}
}

func entityToDomain(sl dao.ShortLink) domain.ShortLink {
	return domain.ShortLink{
		Short: sl.Short,
		Long:  sl.Long,
	}
}
