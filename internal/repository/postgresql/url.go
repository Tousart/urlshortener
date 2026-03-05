package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tousart/urlshortener/internal/domain"
	"gorm.io/gorm"
)

// PostgreSQLRepository - для работы с хранилищем, основанным на PostgreSQL
type PostgreSQLRepository struct {
	db *gorm.DB
}

func NewPostgreSQLRepository(db *gorm.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

// Save() - сохранение пары URL-ов (укороченного и оригинального)
func (r *PostgreSQLRepository) Save(ctx context.Context, url *domain.URL) error {
	err := r.db.WithContext(ctx).Table("urls").Create(url).Error
	if err != nil {
		errPath := "repository: postgresql: Save:"
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			errMsg := err.Error()
			switch {
			case strings.Contains(errMsg, "short"):
				return fmt.Errorf("%s %w", errPath, domain.ErrShortURLExists)
			case strings.Contains(errMsg, "original"):
				return fmt.Errorf("%s %w", errPath, domain.ErrOriginalURLExists)
			}
		}
		return fmt.Errorf("%s %w", errPath, err)
	}
	return nil
}

// GetOriginal() - получение оригинального URL по укороченнному
func (r *PostgreSQLRepository) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	res := r.db.WithContext(ctx).
		Table("urls").
		Select("original_url").
		Where("short_url = ?", shortURL).
		Scan(&originalURL)
	if res.Error != nil {
		return "", fmt.Errorf("repository: postgresql: GetOriginal: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return "", fmt.Errorf("repository: postgresql: GetOriginal: %w", domain.ErrOriginalURLNotFound)
	}
	return originalURL, nil
}

// GetShort() - получение укороченного URL по оригинальному
func (r *PostgreSQLRepository) GetShort(ctx context.Context, originalURL string) (string, error) {
	var shortURL string
	res := r.db.WithContext(ctx).
		Table("urls").
		Select("short_url").
		Where("original_url = ?", originalURL).
		Scan(&shortURL)
	if res.Error != nil {
		return "", fmt.Errorf("repository: postgresql: GetShort: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return "", fmt.Errorf("repository: postgresql: GetShort: %w", domain.ErrShortURLNotFound)
	}
	return shortURL, nil
}
