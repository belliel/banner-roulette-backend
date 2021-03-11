package app

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/apiserver"
	"github.com/BellZaph/banner-roulette-backend/internal/config"
	"github.com/BellZaph/banner-roulette-backend/internal/repository"
	"github.com/BellZaph/banner-roulette-backend/internal/service"
	"github.com/BellZaph/banner-roulette-backend/internal/transport/http"
	"github.com/BellZaph/banner-roulette-backend/pkg/database/mongodb"
	"github.com/BellZaph/banner-roulette-backend/pkg/hash"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Start(configPath string) {
	cfg := config.Get(configPath)

	db, err := mongodb.NewMongoDBDatabase(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	hasher := hash.NewRandHash()

	repositories, err := repository.NewRepository(db)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	services := service.NewServices(service.Deps{
		Repos:  repositories,
		Hasher: hasher,
	})


	handlers := http.NewHandler(services)

	srv := apiserver.NewServer(cfg, handlers.Init("", cfg.HTTPPort))

	go func() {
		if err := srv.Run(); err != nil {
			logrus.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logrus.Info("Server started")
	logrus.Infof("On http://localhost:%s", cfg.HTTPPort)
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := db.Client().Disconnect(context.Background()); err != nil {
		logrus.Error(err.Error())
	}
}

