package service

import "awesome-shortLink/domain"

type ShortLinkService interface {
	ShortenURL(url string) (domain.ShorLink, error)
}

type ShortLinkServiceBasic struct {
	repo ShortLinkRepository
}

func NewShortLinkServiceBasic(repo ShortLinkRepository) ShortLinkService {
	return &ShortLinkServiceBasic{
		repo: repo,
	}
}
