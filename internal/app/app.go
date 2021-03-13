package app

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/apiserver"
	"github.com/BellZaph/banner-roulette-backend/internal/config"
	"github.com/BellZaph/banner-roulette-backend/internal/repository"
	"github.com/BellZaph/banner-roulette-backend/internal/service"
	"github.com/BellZaph/banner-roulette-backend/pkg/database/mongodb"
	"github.com/BellZaph/banner-roulette-backend/pkg/hash"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// @title Banner roulette backend
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name Ruslan
// @contact.email rkuserbaev@gmail.com

// @host localhost:5000
// @BasePath /v1
// @query.collection.format multi
func Start(configPath string) {
	cfg := config.Get(configPath)

	db, err := mongodb.NewMongoDBDatabase(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		logrus.Fatal(err.Error())
	} else {
		logrus.Infof("%s connected", db.Name())
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

	appCtx, cancelAppCtx := context.WithCancel(context.Background())
	defer cancelAppCtx()


	srv := apiserver.NewHTTPServer(appCtx, cfg, services, hasher,"")

	go func() {
		if err := srv.Run(); err != nil {
			logrus.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Warn("interrupt signal")
	cancelAppCtx()

	if err := db.Client().Disconnect(context.Background()); err != nil {
		logrus.Error(err.Error())
	}
}

