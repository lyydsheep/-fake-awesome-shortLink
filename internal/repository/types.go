package repository

import (
	"awesome-shortLink/internal/domain"
	"awesome-shortLink/internal/repository/cache"
	"awesome-shortLink/internal/repository/dao"
	"context"
)

var (
	ErrDuplicatedLongURL = dao.ErrDuplicatedLongURL
	ErrNotFound          = dao.ErrNotFound
	ErrNotExistsInFilter = cache.ErrNotExistsInFilter
)

const (
	bloomFilterKey = "bloomFilter:shortURL"
	countKey       = "CntIdKey"
)

type ShortLinkRepository interface {
	Create(ctx context.Context, longURL string) (domain.ShortLink, error)
	FindByShort(ctx context.Context, shortURL string) (domain.ShortLink, error)
}
