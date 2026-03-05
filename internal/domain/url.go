package domain

import (
	"fmt"
	"net/url"
	"regexp"
)

// URL - доменная модель, которая включает в себя оригинальный URL и его укороченную версию
type URL struct {
	Original string `gorm:"column:original_url"`
	Short    string `gorm:"column:short_url"`
}

// LengthShortURL - длина укороченной ссылки
const LengthShortURL = 10

var shortRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ValidateOriginal() - валидация оригинального URL
func (ur *URL) ValidateOriginal() error {
	const errPath = "domain: ValidateOriginal:"
	if len(ur.Original) < 10 {
		return ErrIncorrectOriginalLength
	}
	parsed, err := url.ParseRequestURI(ur.Original)
	if err != nil {
		return fmt.Errorf("%s %w", errPath, ErrInvalidOriginalURL)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("%s %w", errPath, ErrInvalidOriginalURL)
	}
	return nil
}

// ValidateShort() - валидация укороченного URL
func (ur *URL) ValidateShort() error {
	const errPath = "domain: ValidateShort:"
	if len(ur.Short) != LengthShortURL {
		return fmt.Errorf("%s %w", errPath, ErrIncorrectShortLength)
	}
	if !shortRegex.MatchString(ur.Short) {
		return fmt.Errorf("%s %w", errPath, ErrInvalidShortURL)
	}
	return nil
}
