package main

import (
	"os"

	"github.com/theskyinflames/cdmon2/app/store"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/theskyinflames/cdmon2/app/api"
	"github.com/theskyinflames/cdmon2/app/config"
	"github.com/theskyinflames/cdmon2/app/domain"
	"github.com/theskyinflames/cdmon2/app/repository"
	"github.com/theskyinflames/cdmon2/app/service"
)

func main() {

	// Set logging service
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	log := logrus.New()

	// Sever config loading
	cfg := &config.Config{}
	cfg.Load()
	log.Info(spew.Sdump(cfg))

	// Init the hostings server domain
	serverDomain, err := domain.NewServer(cfg)
	if err != nil {
		panic(err)
	}

	// Init the hostings repository
	store, err := store.NewStore(cfg, log)
	if err != nil {
		panic(err)
	}
	store.Flush() // Empty for each execution.
	hostingsRepository := repository.NewHostingReposytoryMap(cfg, store)

	// Init the hostings server service
	service := service.NewServer(hostingsRepository, serverDomain, cfg, log)

	// Init the controller
	controller := api.NewController(service, log)

	// Start the API
	api := api.NewApi(controller, log, cfg)
	api.Start()
}
