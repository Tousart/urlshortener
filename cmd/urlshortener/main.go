package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/urlshortener/configs"
	"github.com/tousart/urlshortener/internal/api"
	"github.com/tousart/urlshortener/internal/app"
	"github.com/tousart/urlshortener/internal/server"
	"github.com/tousart/urlshortener/internal/usecase"
	"github.com/tousart/urlshortener/pkg/generator"
	"golang.org/x/sync/errgroup"
)

func main() {
	// контекст, ожидающий конкретных системных сигналов (для завершения программы)
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// error-группа для ожидания завершения всех ключевых горутин
	ewg, ctx := errgroup.WithContext(sigCtx)

	// configs
	flags := configs.ParseFlags()
	cfg, err := configs.LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// logger
	logger := app.InitLogger()

	// repository
	urlRepo, err := app.InitRepo(flags.Repository, cfg)
	if err != nil {
		log.Fatalf("failed to init repository: %v", err)
	}

	// generator для генерации укороченного URL
	urlGenerator := generator.NewGenerator()

	// service
	urlUsecase := usecase.NewURLUsecase(urlRepo, urlGenerator)

	// api + router
	r := chi.NewRouter()
	urlAPI := api.NewAPI(urlUsecase, logger)
	urlAPI.WithHandlers(r)

	// server
	srv := server.NewServer(&cfg.Server, r)
	// запуск сервера в отдельной горутине
	ewg.Go(func() error {
		return srv.CreateAndRunServer(ctx)
	})
	// запуск горутины, ожидающей команды на graceful shutdown
	ewg.Go(func() error {
		return srv.ShutdownServer(ctx)
	})

	// ожидание ошибки или отмены контекста в одной из ключевых горутин
	if err := ewg.Wait(); err != nil {
		log.Printf("main error: %v\n", err)
		return
	}
}
