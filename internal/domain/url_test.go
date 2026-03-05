package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// полностью корректные данные
func TestURLValidateSuccess(t *testing.T) {
	url := URL{
		Original: "https://site.com",
		Short:    "aBcDefG1_3",
	}

	err := url.ValidateOriginal()
	assert.NoError(t, err)

	err = url.ValidateShort()
	assert.NoError(t, err)
}

// невалидный формат оригинальной ссылки
func TestURLValidateInvalidOriginal(t *testing.T) {
	url := URL{
		Original: "prosto text ne ssilka",
		Short:    "1234567890",
	}

	err := url.ValidateOriginal()
	assert.ErrorIs(t, err, ErrInvalidOriginalURL)
}

// невалидный формат короткой ссылки
func TestURValidateInvalidShort(t *testing.T) {
	url := URL{
		Original: "https://google.com",
		Short:    "short@link", // @ запрещен
	}

	err := url.ValidateShort()
	assert.ErrorIs(t, err, ErrInvalidShortURL)
}

// неправильная длина короткой ссылки
func TestURLValidateShortIncorrectLength(t *testing.T) {
	url := URL{
		Original: "https://site.com",
		Short:    "too_short", // 9 символов вместо 10
	}

	err := url.ValidateShort()
	assert.ErrorIs(t, err, ErrIncorrectShortLength)
}
