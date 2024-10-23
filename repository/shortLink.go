package repository

import (
	"awesome-shortLink/dao"
	"awesome-shortLink/domain"
	"context"
	"errors"
)

var (
	ErrDuplicatedLongURL = dao.ErrDuplicatedLongURL
)

type ShortLinkRepository interface {
	Create(ctx context.Context, url string) (domain.ShortLink, error)
}

type ShortLinkRepositoryV1 struct {
	dao dao.ShortLinkDAO
}

func (repo *ShortLinkRepositoryV1) Create(ctx context.Context, url string) (domain.ShortLink, error) {
	sl, err := repo.dao.Insert(ctx, url)
	switch {
	case err == nil:
		return entityToDomain(sl), nil
	case errors.Is(err, ErrDuplicatedLongURL):
		return domain.ShortLink{}, ErrDuplicatedLongURL
	default:
		return domain.ShortLink{}, err
	}
}

func NewShortLinkRepositoryV1() ShortLinkRepository {
	return &ShortLinkRepositoryV1{}
}

func entityToDomain(sl dao.ShortLink) domain.ShortLink {
	return domain.ShortLink{
		Short: sl.Short,
		Long:  sl.Long,
	}
}
