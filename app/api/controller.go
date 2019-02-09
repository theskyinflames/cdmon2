package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hako/durafmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/theskyinflames/cdmon2/app"
	"github.com/theskyinflames/cdmon2/app/domain"
)

type (
	ServerService interface {
		CreateHosting(name string, cores int, memorymb int, diskmb int) (domain.UUID, error)
		GetHostings() ([]domain.Hosting, error)
		RemoveHosting(uuid domain.UUID) error
		UpdateHosting(hosting *domain.Hosting) error
		GetServerStatus() *domain.Server
	}

	HealthRs struct {
		RunningTime  string `json:"running_time"`
		ServerStatus *domain.Server
	}

	CreateHostingRq struct {
		domain.Hosting
	}
	CreateHostingRs struct {
		UUID   string `json:"uuid,omitempty"`
		ErrMsg string `json:"error,omitempty"`
	}

	GetHostingsRs struct {
		Hostings []domain.Hosting
		ErrMsg   string `json:"error,omitempty"`
	}

	RemoveHostingRs struct {
		UUID   string `json:"uuid,omitempty"`
		ErrMsg string `json:"error,omitempty"`
	}

	UpdateHostingRq struct {
		domain.Hosting
	}
	UpdateHostingRs struct {
		UUID   string `json:"uuid,omitempty"`
		ErrMsg string `json:"error,omitempty"`
	}

	Controller struct {
		log           *logrus.Logger
		startTime     time.Time
		serverService ServerService
	}
)

func NewController(serverService ServerService, log *logrus.Logger) *Controller {
	return &Controller{
		log:           log,
		serverService: serverService,
		startTime:     time.Now(),
	}
}

func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	var (
		rs HealthRs
	)

	runningTime := time.Now().Sub(c.startTime)
	serverStatus := c.serverService.GetServerStatus()
	rs = HealthRs{RunningTime: durafmt.Parse(runningTime).String(), ServerStatus: serverStatus}
	c.respondWithJson(w, http.StatusOK, &rs, r.Method)
}

func (c *Controller) CreateHosting(w http.ResponseWriter, r *http.Request) {
	var (
		rq CreateHostingRq
		rs CreateHostingRs
	)

	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		c.respondWithJson(w, http.StatusBadRequest, err.Error(), r.Method)
		return
	}

	uuid, err := c.serverService.CreateHosting(rq.Name, rq.Cores, rq.MemoryMb, rq.DiskMb)
	if err != nil {
		rs = CreateHostingRs{ErrMsg: err.Error()}
		switch errors.Cause(err) {
		case app.DbErrorAlreadyExist:
			c.respondWithJson(w, http.StatusConflict, &rs, r.Method)
		default:
			c.respondWithJson(w, http.StatusInternalServerError, &rs, r.Method)
		}
		return
	}

	rs = CreateHostingRs{UUID: string(uuid)}
	c.respondWithJson(w, http.StatusOK, &rs, r.Method)
}

func (c *Controller) GetHostings(w http.ResponseWriter, r *http.Request) {
	var (
		rs GetHostingsRs
	)

	hostings, err := c.serverService.GetHostings()
	if err != nil {
		rs = GetHostingsRs{ErrMsg: err.Error()}
		c.respondWithJson(w, http.StatusInternalServerError, &rs, r.Method)
		return
	}

	rs = GetHostingsRs{Hostings: hostings}
	if len(hostings) > 0 {
		c.respondWithJson(w, http.StatusFound, &rs, r.Method)
	} else {
		c.respondWithJson(w, http.StatusOK, &rs, r.Method)
	}
}

func (c *Controller) RemoveHosting(w http.ResponseWriter, r *http.Request) {
	var rs RemoveHostingRs

	params := mux.Vars(r)
	uuid := params["uuid"]

	err := c.serverService.RemoveHosting(domain.UUID(uuid))
	if err != nil {
		rs = RemoveHostingRs{ErrMsg: err.Error()}
		switch errors.Cause(err) {
		case app.DbErrorNotFound:
			c.respondWithJson(w, http.StatusNotFound, &rs, r.Method)
		default:
			c.respondWithJson(w, http.StatusInternalServerError, &rs, r.Method)
		}
		return
	}

	rs = RemoveHostingRs{UUID: uuid}
	c.respondWithJson(w, http.StatusOK, rs, r.Method)
}

func (c *Controller) UpdateHosting(w http.ResponseWriter, r *http.Request) {
	var (
		rq UpdateHostingRq
		rs UpdateHostingRs
	)

	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		c.respondWithJson(w, http.StatusBadRequest, err.Error(), r.Method)
		return
	}

	err = c.serverService.UpdateHosting(&rq.Hosting)
	if err != nil {
		rs = UpdateHostingRs{ErrMsg: err.Error()}
		switch errors.Cause(err) {
		case app.DbErrorNotFound:
			c.respondWithJson(w, http.StatusNotFound, &rs, r.Method)
		default:
			c.respondWithJson(w, http.StatusInternalServerError, &rs, r.Method)
		}
		return
	}

	rs = UpdateHostingRs{UUID: string(rq.UUID)}
	c.respondWithJson(w, http.StatusOK, &rs, r.Method)
}

func (c *Controller) respondWithJson(w http.ResponseWriter, code int, payload interface{}, action string) {
	response, _ := json.Marshal(payload)

	if code != http.StatusOK && code != http.StatusFound && response != nil {
		c.log.WithFields(logrus.Fields{
			"action":      action,
			"http_status": code,
		}).Error(string(response))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
