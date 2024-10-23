package service

import (
	"awesome-shortLink/domain"
	"awesome-shortLink/repository"
	"context"
)

type ShortLinkService interface {
	ShortenURL(ctx context.Context, url string) (domain.ShortLink, error)
}

type ShortLinkServiceBasic struct {
	repo repository.ShortLinkRepository
}

func (svc *ShortLinkServiceBasic) ShortenURL(ctx context.Context, url string) (domain.ShortLink, error) {
	return svc.repo.Create(ctx, url)
}

func NewShortLinkServiceBasic(repo repository.ShortLinkRepository) ShortLinkService {
	return &ShortLinkServiceBasic{
		repo: repo,
	}
}
