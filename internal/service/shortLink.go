package service

import (
	"awesome-shortLink/internal/domain"
	"awesome-shortLink/internal/repository"
	"context"
	"go.uber.org/zap"
)

type ShortLinkService interface {
	ShortenURL(ctx context.Context, longURL string) (domain.ShortLink, error)
	Obtain(ctx context.Context, shortURL string) (domain.ShortLink, error)
}

type ShortLinkServiceBasic struct {
	repo repository.ShortLinkRepository
	l    *zap.Logger
}

func (svc *ShortLinkServiceBasic) Obtain(ctx context.Context, shortURL string) (domain.ShortLink, error) {
	return svc.repo.FindByShort(ctx, shortURL)
}

func (svc *ShortLinkServiceBasic) ShortenURL(ctx context.Context, longURL string) (domain.ShortLink, error) {
	return svc.repo.Create(ctx, longURL)
}

func NewShortLinkServiceBasic(repo repository.ShortLinkRepository, l *zap.Logger) ShortLinkService {
	return &ShortLinkServiceBasic{
		repo: repo,
		l:    l,
	}
}
