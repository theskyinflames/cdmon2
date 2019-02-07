package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hako/durafmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	service "github.com/theskyinflames/cdmon2/hosting"

	"github.com/theskyinflames/cdmon2/hosting/domain"
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
	respondWithJson(w, http.StatusOK, &rs)
}

func (c *Controller) CreateHosting(w http.ResponseWriter, r *http.Request) {
	var (
		rq CreateHostingRq
		rs CreateHostingRs
	)

	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, err.Error())
		return
	}

	uuid, err := c.serverService.CreateHosting(rq.Name, rq.Cores, rq.MemoryMb, rq.DiskMb)
	if err != nil {
		rs = CreateHostingRs{ErrMsg: err.Error()}
		switch errors.Cause(err).(type) {
		case service.DbErrorAlreadyExist:
			respondWithJson(w, http.StatusConflict, &rs)
		default:
			respondWithJson(w, http.StatusInternalServerError, &rs)
		}
	}

	rs = CreateHostingRs{UUID: string(uuid)}
	respondWithJson(w, http.StatusOK, &rs)
}

func (c *Controller) GetHostings(w http.ResponseWriter, r *http.Request) {
	var (
		rs GetHostingsRs
	)

	hostings, err := c.serverService.GetHostings()
	if err != nil {
		rs = GetHostingsRs{ErrMsg: err.Error()}
		respondWithJson(w, http.StatusInternalServerError, &rs)
		return
	}

	rs = GetHostingsRs{Hostings: hostings}
	if len(hostings) > 0 {
		respondWithJson(w, http.StatusFound, &rs)
	} else {
		respondWithJson(w, http.StatusOK, &rs)
	}
}

func (c *Controller) RemoveHosting(w http.ResponseWriter, r *http.Request) {
	var rs RemoveHostingRs

	params := mux.Vars(r)
	uuid := params["uuid"]

	err := c.serverService.RemoveHosting(domain.UUID(uuid))
	if err != nil {
		rs = RemoveHostingRs{ErrMsg: err.Error()}
		switch errors.Cause(err).(type) {
		case service.DbErrorNotFound:
			respondWithJson(w, http.StatusNotFound, &rs)
		default:
			respondWithJson(w, http.StatusInternalServerError, &rs)
		}
		return
	}

	rs = RemoveHostingRs{UUID: uuid}
	respondWithJson(w, http.StatusOK, rs)
}

func (c *Controller) UpdateHosting(w http.ResponseWriter, r *http.Request) {
	var (
		rq UpdateHostingRq
		rs UpdateHostingRs
	)

	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.serverService.UpdateHosting(&rq.Hosting)
	if err != nil {
		rs = UpdateHostingRs{ErrMsg: err.Error()}
		switch errors.Cause(err).(type) {
		case service.DbErrorNotFound:
			respondWithJson(w, http.StatusNotFound, &rs)
		default:
			respondWithJson(w, http.StatusInternalServerError, &rs)
		}
		return
	}

	rs = UpdateHostingRs{UUID: string(rq.UUID)}
	respondWithJson(w, http.StatusOK, &rs)

}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
