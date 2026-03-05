package app

import (
	"fmt"

	"github.com/tousart/urlshortener/configs"
	"github.com/tousart/urlshortener/internal/domain"
	"github.com/tousart/urlshortener/internal/repository/inmemory"
	"github.com/tousart/urlshortener/internal/repository/postgresql"
	"github.com/tousart/urlshortener/internal/usecase"
	pkgpsql "github.com/tousart/urlshortener/pkg/postgresql"
)

// InitRepo() - инициализация репозитория в зависимости от типа (переданного через флаг -repo)
func InitRepo(dbType string, cfg *configs.Config) (usecase.URLRepository, error) {
	var urlRepo usecase.URLRepository
	switch dbType {
	case "inmemory":
		urlRepo = inmemory.NewInMemoryRepository()
	case "postgresql":
		db, err := pkgpsql.ConnectToPSQL(&cfg.PostgreSQL)
		if err != nil {
			return nil, fmt.Errorf("app: InitRepo: %w", err)
		}
		urlRepo = postgresql.NewPostgreSQLRepository(db)
	default:
		return nil, domain.ErrUnknownRepoType
	}
	return urlRepo, nil
}
