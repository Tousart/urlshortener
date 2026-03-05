package usecase

import (
	"context"

	"github.com/tousart/urlshortener/internal/domain"
)

type URLRepository interface {
	Save(ctx context.Context, url *domain.URL) error
	GetOriginal(ctx context.Context, shortURL string) (string, error)
	GetShort(ctx context.Context, originalURL string) (string, error)
}

type URLGenerator interface {
	GenerateRandomStringWithLength(length int) (string, error)
}
