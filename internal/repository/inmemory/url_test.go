package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tousart/urlshortener/internal/domain"
)

/*
	Save()
*/

func TestInMemoryRepositorySaveSuccess(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()
	url := &domain.URL{Original: "https://site.com", Short: "abcdxyz_21"}

	err := repo.Save(ctx, url)

	assert.NoError(t, err)

	savedOrig, _ := repo.GetOriginal(ctx, "abcdxyz_21")
	assert.Equal(t, url.Original, savedOrig)

	savedShort, _ := repo.GetShort(ctx, "https://site.com")
	assert.Equal(t, url.Short, savedShort)
}

func TestInMemoryRepositorySaveDuplicateShort(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()
	url1 := &domain.URL{Original: "https://ozonbank.ru", Short: "duplicate0"}
	url2 := &domain.URL{Original: "https://ozonfresh.ru", Short: "duplicate0"}

	_ = repo.Save(ctx, url1)

	err := repo.Save(ctx, url2)
	assert.ErrorIs(t, err, domain.ErrShortURLExists)
}

func TestInMemoryRepositorySaveDuplicateOriginal(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()
	url1 := &domain.URL{Original: "https://ozon.ru", Short: "duplicate0"}
	url2 := &domain.URL{Original: "https://ozon.ru", Short: "duplicate1"}

	_ = repo.Save(ctx, url1)

	err := repo.Save(ctx, url2)
	assert.ErrorIs(t, err, domain.ErrOriginalURLExists)
}

/*
	GetOriginal()
*/

func TestInMemoryRepositoryGetOriginalSuccess(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()
	originalURL := "https://site.com"
	shortURL := "abcdxyz_21"

	_ = repo.Save(ctx, &domain.URL{Original: originalURL, Short: shortURL})

	actual, err := repo.GetOriginal(ctx, shortURL)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, actual)
}

func TestInMemoryRepositoryGetOriginalNotFound(t *testing.T) {
	repo := NewInMemoryRepository()

	actual, err := repo.GetOriginal(context.Background(), "abcdxyz_21")

	assert.Empty(t, actual)
	assert.ErrorIs(t, err, domain.ErrOriginalURLNotFound)
}

/*
	GetShort()
*/

func TestInMemoryRepositoryGetShortSuccess(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()
	originalURL := "https://site.com"
	shortURL := "abcdxyz_21"

	_ = repo.Save(ctx, &domain.URL{Original: originalURL, Short: shortURL})

	actual, err := repo.GetShort(ctx, originalURL)

	assert.NoError(t, err)
	assert.Equal(t, shortURL, actual)
}

func TestInMemoryRepositoryGetShortNotFound(t *testing.T) {
	repo := NewInMemoryRepository()

	actual, err := repo.GetShort(context.Background(), "https://site.com")

	assert.Empty(t, actual)
	assert.ErrorIs(t, err, domain.ErrShortURLNotFound)
}
