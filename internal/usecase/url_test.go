package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tousart/urlshortener/internal/domain"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetShort(ctx context.Context, originalURL string) (string, error) {
	args := m.Called(ctx, originalURL)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	args := m.Called(ctx, shortURL)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) Save(ctx context.Context, url *domain.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

type mockGen struct {
	mock.Mock
}

func (m *mockGen) GenerateRandomStringWithLength(length int) (string, error) {
	args := m.Called(length)
	return args.String(0), args.Error(1)
}

/*
	Shorten()
*/

func TestShortenSuccess(t *testing.T) {
	repo := &mockRepo{}
	gen := &mockGen{}
	urlUsecase := NewURLUsecase(repo, gen)
	originalURL, shortURL := "https://site.com", "aBcDefG1_3"

	// в БД нет такого URL (это нормальное поведение)
	repo.On("GetShort", mock.Anything, originalURL).Return("", domain.ErrShortURLNotFound)

	gen.On("GenerateRandomStringWithLength", domain.LengthShortURL).Return(shortURL, nil)

	// проверяем, что в repository отправилась корректная доменная модель
	repo.On("Save", mock.Anything, mock.MatchedBy(func(url *domain.URL) bool {
		return url.Original == originalURL && url.Short == shortURL
	})).Return(nil)

	actualShortURL, err := urlUsecase.Shorten(context.Background(), originalURL)

	assert.NoError(t, err)
	assert.Equal(t, shortURL, actualShortURL)
}

// такой URL уже сокращали, поэтому нужно просто вернуть укороченный URL (не делать укорачивание заново)
func TestShortenSuccessAlreadyExists(t *testing.T) {
	repo := &mockRepo{}
	gen := &mockGen{}
	urlUsecase := NewURLUsecase(repo, gen)
	originalURL, shortURL := "http://localhost:6666", "aBcDefG1_3"

	repo.On("GetShort", mock.Anything, originalURL).Return(shortURL, nil)

	actualShortURL, err := urlUsecase.Shorten(context.Background(), originalURL)

	assert.NoError(t, err)
	assert.Equal(t, shortURL, actualShortURL)
	gen.AssertNotCalled(t, "GenerateRandomStringWithLength", mock.Anything)
}

// передаем невалидный URL
func TestShortenErrorInvalidURL(t *testing.T) {
	urlUsecase := NewURLUsecase(new(mockRepo), new(mockGen))

	actualShortURL, err := urlUsecase.Shorten(context.Background(), "stroka-normalnoy-dlini")

	assert.Error(t, err)
	assert.Empty(t, actualShortURL)
	assert.Contains(t, err.Error(), "invalid original url format")
}

/*
	GetOriginal()
*/

// успешный возврат оригинального URL по укороченному
func TestGetOriginalSuccess(t *testing.T) {
	repo := &mockRepo{}
	urlUsecase := NewURLUsecase(repo, nil)
	shortURL, originalURL := "aBcDefG1_3", "https://saytyan.net"

	repo.On("GetOriginal", mock.Anything, shortURL).Return(originalURL, nil)

	actualOriginalURL, err := urlUsecase.GetOriginal(context.Background(), shortURL)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, actualOriginalURL)
}

// для такого укороченного URL не найдено оригинального
func TestGetOriginalNotFound(t *testing.T) {
	repo := &mockRepo{}
	urlUsecase := NewURLUsecase(repo, nil)
	shortURL := "aBcDefG1_3"

	repo.On("GetOriginal", mock.Anything, shortURL).Return("", domain.ErrOriginalURLNotFound)

	actualOriginalURL, err := urlUsecase.GetOriginal(context.Background(), shortURL)

	assert.ErrorIs(t, err, domain.ErrOriginalURLNotFound)
	assert.Empty(t, actualOriginalURL)
}

// передан невалидный укороченный URL
func TestGetOriginal(t *testing.T) {
	repo := &mockRepo{}
	urlUsecase := NewURLUsecase(repo, nil)
	shortURL := "aBcD*??:hg" // *, ? и : - запрещенные символы

	actualOriginalURL, err := urlUsecase.GetOriginal(context.Background(), shortURL)

	assert.ErrorIs(t, err, domain.ErrInvalidShortURL)
	assert.Empty(t, actualOriginalURL)
	repo.AssertNotCalled(t, "GetOriginal", mock.Anything)
}
