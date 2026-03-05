package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/tousart/urlshortener/internal/domain"
)

// URLUsecase - для реализации бизнес-логики с укорачиванием ссылок
type URLUsecase struct {
	urlRepo      URLRepository
	urlGenerator URLGenerator
}

func NewURLUsecase(urlRepo URLRepository, urlGenerator URLGenerator) *URLUsecase {
	return &URLUsecase{
		urlRepo:      urlRepo,
		urlGenerator: urlGenerator,
	}
}

// Shorten() - укорачивание оригинального URL
func (u *URLUsecase) Shorten(ctx context.Context, originalURL string) (string, error) {
	const errPath = "usecase: Shorten error:"

	// создаем доменную модель URL и валидируем оригинальный URL
	url := &domain.URL{Original: originalURL}
	if err := url.ValidateOriginal(); err != nil {
		return "", fmt.Errorf("%s %w", errPath, err)
	}

	// проверка: существует ли укороченный URL по оригинальному из запроса?
	shortURL, err := u.urlRepo.GetShort(ctx, originalURL)
	if err == nil {
		return shortURL, nil
	} else if !errors.Is(err, domain.ErrShortURLNotFound) {
		return "", fmt.Errorf("%s %w", errPath, err)
	}

	// укорачиваем URL
	shortURL, err = u.urlGenerator.GenerateRandomStringWithLength(domain.LengthShortURL)
	if err != nil {
		return "", fmt.Errorf("%s %w", errPath, err)
	}

	// валидируем укороченный URL
	url.Short = shortURL
	if err = url.ValidateShort(); err != nil {
		return "", fmt.Errorf("%s %w", errPath, err)
	}

	// сохраняем URL-ы в БД
	if err = u.urlRepo.Save(ctx, url); err != nil {
		// если кто-то конкурентно быстрее сохранил оригинальный URL из запроса, то получаем укороченный
		if errors.Is(err, domain.ErrOriginalURLExists) {
			url.Short, err = u.urlRepo.GetShort(ctx, originalURL)
			if err != nil {
				return "", fmt.Errorf("%s, %w", errPath, err)
			}
		}
		return "", fmt.Errorf("%s %w", errPath, err)
	}
	return shortURL, nil
}

// GetOriginal() - получение оригинального URL по укороченному
func (u *URLUsecase) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	url := domain.URL{Short: shortURL}
	if err := url.ValidateShort(); err != nil {
		return "", fmt.Errorf("usecase: GetOriginal error: %w", err)
	}

	originalURL, err := u.urlRepo.GetOriginal(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("usecase: GetOriginal error: %w", err)
	}
	return originalURL, nil
}
