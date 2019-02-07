package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/theskyinflames/cdmon2/hosting/api"
	"github.com/theskyinflames/cdmon2/hosting/application"
	"github.com/theskyinflames/cdmon2/hosting/config"
	"github.com/theskyinflames/cdmon2/hosting/domain"
	"github.com/theskyinflames/cdmon2/hosting/repository"
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
	hostingsRepository := repository.NewHostingReposytoryMap(cfg)

	// Init the hostings server service
	service := application.NewServer(hostingsRepository, serverDomain, cfg, log)

	// Init the controller
	controller := api.NewController(service, log)

	// Start the API
	api := api.NewApi(controller, log, cfg)
	api.Start()
}
