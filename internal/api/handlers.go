package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/urlshortener/internal/domain"
	"github.com/tousart/urlshortener/internal/middleware"
)

// URLUsecase - интерфейс к сервису по укорачиванию URL-ов
type URLUsecase interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	GetOriginal(ctx context.Context, shortURL string) (string, error)
}

// API - для обработки запросов по HTTP
type API struct {
	urlUsecase URLUsecase
	logger     *slog.Logger
}

func NewAPI(urlUsecase URLUsecase, logger *slog.Logger) *API {
	return &API{
		urlUsecase: urlUsecase,
		logger:     logger,
	}
}

// ShortenURLRequest - запрос на сокращение URL
type ShortenURLRequest struct {
	URL string `json:"url"`
}

// ShortenURLResponse - ответ на сокращение URL
type ShortenURLResponse struct {
	ShortURL string `json:"short_url"`
}

// GetOriginalURLResponse - ответ на получение оригинального URL по укороченному
type GetOriginalURLResponse struct {
	URL string `json:"url"`
}

// ErrorResponse - ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// WithHandlers() - назначаем хэндлеры
func (ap *API) WithHandlers(r chi.Router) {
	r.Use(middleware.LoggingMiddleware(ap.logger))

	r.Route("/urls", func(r chi.Router) {
		r.Post("/", ap.shortenURLHandler)
		r.Get("/{short}", ap.getOriginalURLHandler)
	})
}

// обработчик укорачивания URL
func (ap *API) shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ap.logger.Info("bad request")
		ap.responseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		ap.logger.Info("bad request")
		ap.responseError(w, http.StatusBadRequest, "url is required")
		return
	}

	short, err := ap.urlUsecase.Shorten(r.Context(), req.URL)
	if err != nil {
		respStatus, errMsg := ap.processError(err, slog.String("url", req.URL))
		ap.responseError(w, respStatus, errMsg)
		return
	}

	ap.responseJSON(w, http.StatusCreated, ShortenURLResponse{ShortURL: short})
}

// обработчик получения оригинального URL по укороченному
func (ap *API) getOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "short")
	if shortURL == "" {
		ap.logger.Info("bad request")
		ap.responseError(w, http.StatusBadRequest, "short url is required")
		return
	}

	originalURL, err := ap.urlUsecase.GetOriginal(r.Context(), shortURL)
	if err != nil {
		respStatus, errMsg := ap.processError(err, slog.String("short_url", shortURL))
		ap.responseError(w, respStatus, errMsg)
		return
	}

	ap.responseJSON(w, http.StatusOK, GetOriginalURLResponse{URL: originalURL})
}

func (ap *API) responseJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if response != nil {
		json.NewEncoder(w).Encode(response)
	}
}

func (ap *API) responseError(w http.ResponseWriter, status int, errMsg string) {
	ap.responseJSON(w, status, ErrorResponse{Error: errMsg})
}

// возвращает для ошибки: код, текст
func (ap *API) processError(err error, args ...any) (int, string) {
	var (
		internalErr bool
		respStatus  int
		errMsg      string
	)

	switch {
	case errors.Is(err, domain.ErrIncorrectOriginalLength):
		errMsg = domain.ErrIncorrectOriginalLength.Error()
		respStatus = http.StatusBadRequest

	case errors.Is(err, domain.ErrInvalidOriginalURL):
		errMsg = domain.ErrInvalidOriginalURL.Error()
		respStatus = http.StatusBadRequest

	case errors.Is(err, domain.ErrOriginalURLNotFound):
		errMsg = domain.ErrOriginalURLNotFound.Error()
		respStatus = http.StatusNotFound

	default:
		errMsg = "internal error"
		respStatus = http.StatusInternalServerError
		internalErr = true
	}

	if internalErr {
		ap.logger.Error(errMsg, append(args, slog.Any("error", err))...)
	} else {
		ap.logger.Info(errMsg, append(args, slog.Int("status", respStatus))...)
	}

	return respStatus, errMsg
}
