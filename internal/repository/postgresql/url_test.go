package postgresql

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tousart/urlshortener/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func mockSQL(t *testing.T) (*PostgreSQLRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	require.NoError(t, err)

	return NewPostgreSQLRepository(gormDB), mock
}

/*
	Save()
*/

// успешное сохранение
func TestPostgreSQLRepositorySaveSuccess(t *testing.T) {
	repo, mock := mockSQL(t)
	ctx := context.Background()
	url := &domain.URL{Original: "https://site.com", Short: "abcdxyz_21"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "urls"`)).
		WithArgs(url.Original, url.Short).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(ctx, url)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// попытка сохранить дубликат укороченного URL
func TestPostgreSQLRepositorySaveDuplicateShort(t *testing.T) {
	repo, mock := mockSQL(t)
	ctx := context.Background()
	url := &domain.URL{Original: "https://ozon.ru", Short: "duplicate0"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "urls"`)).
		WillReturnError(fmt.Errorf("%w: short", gorm.ErrDuplicatedKey)) // ошибка нарушения уникальности (содержит слово short)
	mock.ExpectRollback()

	err := repo.Save(ctx, url)

	assert.ErrorIs(t, err, domain.ErrShortURLExists)
}

// попытка сохранить дубликат оригинального URL
func TestPostgreSQLRepositorySaveDuplicateOriginal(t *testing.T) {
	repo, mock := mockSQL(t)
	ctx := context.Background()
	url := &domain.URL{Original: "http://localhost:1234", Short: "duplicate1"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "urls"`)).
		WillReturnError(fmt.Errorf("%w: original", gorm.ErrDuplicatedKey)) // ошибка нарушения уникальности (содержит слово original)
	mock.ExpectRollback()

	err := repo.Save(ctx, url)

	assert.ErrorIs(t, err, domain.ErrOriginalURLExists)
}

/*
	GetOriginal()
*/

// успешное возвращение оригинального URL по укороченному
func TestPostgreSQLRepositoryGetOriginalSuccess(t *testing.T) {
	repo, mock := mockSQL(t)
	shortURL := "abcdxyz_21"
	expectedOriginal := "https://site.com"

	rows := sqlmock.NewRows([]string{"original_url"}).AddRow(expectedOriginal)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT original_url FROM "urls"`)).
		WithArgs(shortURL).
		WillReturnRows(rows)

	actual, err := repo.GetOriginal(context.Background(), shortURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedOriginal, actual)
}

// запрос на несуществующий укороченный URL
func TestPostgreSQLRepositoryGetOriginalNotFound(t *testing.T) {
	repo, mock := mockSQL(t)
	ctx := context.Background()
	shortURL := "abcdxyz_21"

	rows := sqlmock.NewRows([]string{"original_url"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT original_url FROM "urls"`)).
		WithArgs(shortURL).
		WillReturnRows(rows)

	actual, err := repo.GetOriginal(ctx, shortURL)

	assert.Empty(t, actual)
	assert.ErrorIs(t, err, domain.ErrOriginalURLNotFound)
}

/*
	GetShort()
*/

// успешное возвращение укороченнного URL по оригинальному
func TestPostgreSQLRepositoryGetShortSuccess(t *testing.T) {
	repo, mock := mockSQL(t)
	ctx := context.Background()
	originalURL := "https://site.com"
	expectedShort := "abcdxyz_21"

	rows := sqlmock.NewRows([]string{"short_url"}).AddRow(expectedShort)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT short_url FROM "urls"`)).
		WithArgs(originalURL).
		WillReturnRows(rows)

	actual, err := repo.GetShort(ctx, originalURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort, actual)
}

// запрос на несуществующий оригинальный URL
func TestPostgreSQLRepositoryGetShortNotFound(t *testing.T) {
	repo, mock := mockSQL(t)
	originalURL := "https://site.com"

	rows := sqlmock.NewRows([]string{"short_url"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT short_url FROM "urls"`)).
		WithArgs(originalURL).
		WillReturnRows(rows)

	actual, err := repo.GetShort(context.Background(), originalURL)

	assert.Empty(t, actual)
	assert.ErrorIs(t, err, domain.ErrShortURLNotFound)
}
