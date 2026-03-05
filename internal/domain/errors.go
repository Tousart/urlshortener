package domain

import "errors"

var (
	ErrUnknownRepoType error = errors.New("unknown type of repository")
	// short url
	ErrShortURLExists       error = errors.New("short URL already exists")
	ErrShortURLNotFound     error = errors.New("short URL not found")
	ErrInvalidShortURL      error = errors.New("short URL must contain only letters (a-z, A-Z), digits, and underscore")
	ErrIncorrectShortLength error = errors.New("length of short URL must be 10 characters")
	// original url
	ErrOriginalURLExists       error = errors.New("original URL already exists")
	ErrOriginalURLNotFound     error = errors.New("original URL not found")
	ErrInvalidOriginalURL      error = errors.New("invalid original url format")
	ErrIncorrectOriginalLength error = errors.New("incorrect original length")
)
