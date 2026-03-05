package inmemory

import (
	"context"
	"sync"

	"github.com/tousart/urlshortener/internal/domain"
)

// InMemoryRepository - для работы с хранилищем, основанным на map
type InMemoryRepository struct {
	shortOriginal map[string]string
	originalShort map[string]string
	mu            *sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		shortOriginal: make(map[string]string),
		originalShort: make(map[string]string),
		mu:            &sync.RWMutex{},
	}
}

// Save() - сохранение пары URL-ов (укороченного и оригинального)
func (r *InMemoryRepository) Save(ctx context.Context, url *domain.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// проверка на наличие сокращенного URL в БД
	if _, ok := r.shortOriginal[url.Short]; ok {
		return domain.ErrShortURLExists
	}

	// проверка на наличие оригинального URL в БД
	if _, ok := r.originalShort[url.Original]; ok {
		return domain.ErrOriginalURLExists
	}

	// если все проверки пройдены, то добавляем ссылку и ее сокращенный код в БД
	r.shortOriginal[url.Short] = url.Original
	r.originalShort[url.Original] = url.Short
	return nil
}

// GetOriginal() - получение оригинального URL по укороченнному
func (r *InMemoryRepository) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	originalURL, ok := r.shortOriginal[shortURL]
	if !ok {
		return "", domain.ErrOriginalURLNotFound
	}
	return originalURL, nil
}

// GetShort() - получение укороченного URL по оригинальному
func (r *InMemoryRepository) GetShort(ctx context.Context, originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shortURL, ok := r.originalShort[originalURL]
	if !ok {
		return "", domain.ErrShortURLNotFound
	}
	return shortURL, nil
}
