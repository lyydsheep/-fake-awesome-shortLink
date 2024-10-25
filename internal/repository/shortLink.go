package repository

import (
	"awesome-shortLink/internal/domain"
	"awesome-shortLink/internal/repository/dao"
	"context"
	"go.uber.org/zap"
)

var (
	ErrDuplicatedLongURL = dao.ErrDuplicatedLongURL
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

func entityToDomain(sl dao.ShortLink) domain.ShortLink {
	return domain.ShortLink{
		Short: sl.Short,
		Long:  sl.Long,
	}
}
