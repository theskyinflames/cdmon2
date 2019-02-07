package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/theskyinflames/cdmon2/hosting/config"
)

type (
	Api struct {
		cfg        *config.Config
		log        *logrus.Logger
		controller *Controller
	}
)

func NewApi(controller *Controller, log *logrus.Logger, cfg *config.Config) *Api {
	return &Api{
		cfg:        cfg,
		log:        log,
		controller: controller,
	}
}

func (a Api) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/health", a.controller.Health).Methods(http.MethodGet)
	router.HandleFunc("/hosting", a.controller.CreateHosting).Methods(http.MethodPost)
	router.HandleFunc("/hosting", a.controller.GetHostings).Methods(http.MethodGet)
	router.HandleFunc("/hosting/{uuid}", a.controller.RemoveHosting).Methods(http.MethodDelete)
	router.HandleFunc("/hosting", a.controller.UpdateHosting).Methods(http.MethodPut)

	a.log.Infof("starting hosting service at port %s", a.cfg.APIPort)
	a.log.Info(http.ListenAndServe(":"+a.cfg.APIPort, router))
}
