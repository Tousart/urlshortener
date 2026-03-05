package postgresql

import (
	"fmt"

	"github.com/tousart/urlshortener/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToPSQL(cfg *configs.PostgreSQLCfg) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("pkg: postgresql: ConnectToPSQL: %w", err)
	}
	return db, nil
}
