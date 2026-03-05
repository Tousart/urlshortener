package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tousart/urlshortener/internal/domain"
)

// мокаем usecase
type mockUsecase struct {
	mock.Mock
}

func (m *mockUsecase) Shorten(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func (m *mockUsecase) GetOriginal(ctx context.Context, short string) (string, error) {
	args := m.Called(ctx, short)
	return args.String(0), args.Error(1)
}

func setupAPI() (*API, *mockUsecase) {
	mockURLUsecase := &mockUsecase{}
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewAPI(mockURLUsecase, discardLogger), mockURLUsecase
}

// POST /urls (укорачивание URL)

// успешный запрос
func TestShortenSuccess(t *testing.T) {
	urlAPI, urlUsecase := setupAPI()

	urlUsecase.On("Shorten", mock.Anything, "https://site.com").Return("aBcDefG1_3", nil)

	body := `{"url": "https://site.com"}`
	req := httptest.NewRequest(http.MethodPost, "/urls", strings.NewReader(body))
	rr := httptest.NewRecorder()

	urlAPI.shortenURLHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, `{"short_url":"aBcDefG1_3"}`, rr.Body.String())
}

// запрос с пустым значением url
func TestShortenEmptyURL(t *testing.T) {
	urlAPI, _ := setupAPI()

	body := `{"url": ""}`
	req := httptest.NewRequest(http.MethodPost, "/urls", strings.NewReader(body))
	rr := httptest.NewRecorder()

	urlAPI.shortenURLHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "url is required")
}

// запрос с невалидным json
func TestShortenInvalidJSON(t *testing.T) {
	urlAPI, _ := setupAPI()

	req := httptest.NewRequest(http.MethodPost, "/urls", strings.NewReader(`{invalid json`))
	rr := httptest.NewRecorder()

	urlAPI.shortenURLHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid request body")
}

// GET /urls/{short} (получение оригинального URL по укороченному)

// успешный запрос
func TestGetOriginalSuccess(t *testing.T) {
	urlAPI, urlUsecase := setupAPI()

	r := chi.NewRouter()
	urlAPI.WithHandlers(r)

	urlUsecase.On("GetOriginal", mock.Anything, "aBcDefG1_3").Return("http://site.com", nil)

	req := httptest.NewRequest(http.MethodGet, "/urls/aBcDefG1_3", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"url":"http://site.com"}`, rr.Body.String())
}

// оригинальная ссылка не найдена
func TestGetOriginalNotFound(t *testing.T) {
	urlAPI, urlUsecase := setupAPI()

	r := chi.NewRouter()
	urlAPI.WithHandlers(r)

	urlUsecase.On("GetOriginal", mock.Anything, "aBcDefG1_3").Return("", domain.ErrOriginalURLNotFound)

	req := httptest.NewRequest(http.MethodGet, "/urls/aBcDefG1_3", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "not found")
}
