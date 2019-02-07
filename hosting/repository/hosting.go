package repository

import (
	"errors"
	"sync"

	service "github.com/theskyinflames/cdmon2/hosting"
	"github.com/theskyinflames/cdmon2/hosting/config"
	"github.com/theskyinflames/cdmon2/hosting/domain"
)

type (
	HostingRepostitoryMap struct {
		sync.Mutex
		cfg       *config.Config
		store     map[domain.UUID]*domain.Hosting
		idxByName map[string]*domain.Hosting
	}
)

func NewHostingReposytoryMap(cfg *config.Config) *HostingRepostitoryMap {
	return &HostingRepostitoryMap{
		Mutex:     sync.Mutex{},
		cfg:       cfg,
		store:     make(map[domain.UUID]*domain.Hosting),
		idxByName: make(map[string]*domain.Hosting),
	}
}

func (h HostingRepostitoryMap) Get(uuid domain.UUID) (*domain.Hosting, error) {
	hosting, exist := h.store[uuid]
	if !exist {
		return nil, service.DbErrorNotFound(errors.New("the hosting doesn't exist"))
	}
	return hosting, nil
}

func (h HostingRepostitoryMap) GetAll() ([]domain.Hosting, error) {
	hostings := make([]domain.Hosting, len(h.store))
	z := 0
	for _, hosting := range h.store {
		hostings[z] = *hosting
		z++
	}
	return hostings, nil
}

func (h HostingRepostitoryMap) Insert(hosting *domain.Hosting) error {
	_, exist := h.store[hosting.UUID]
	if exist {
		return service.DbErrorAlreadyExist(errors.New("the hosting already exist with this uuid"))
	}
	_, exist = h.idxByName[hosting.Name]
	if exist {
		return service.DbErrorAlreadyExist(errors.New("the hosting already exist with this name"))
	}

	err := hosting.Validate(h.cfg)
	if err != nil {
		return err
	}

	h.store[hosting.UUID] = hosting
	h.idxByName[hosting.Name] = hosting
	return nil
}

func (h HostingRepostitoryMap) Update(hosting *domain.Hosting) error {
	old, exist := h.store[hosting.UUID]
	if !exist {
		return service.DbErrorNotFound(errors.New("the hosting doesn't exist"))
	}

	if old.Name != hosting.Name {
		_, exist = h.idxByName[hosting.Name]
		if exist {
			return service.DbErrorAlreadyExist(errors.New("the hosting already exist with this name"))
		}
	}

	err := hosting.Validate(h.cfg)
	if err != nil {
		return err
	}

	h.store[hosting.UUID] = hosting
	h.idxByName[hosting.Name] = hosting
	return nil
}

func (h HostingRepostitoryMap) Remove(uuid domain.UUID) (*domain.Hosting, error) {
	hosting, exist := h.store[uuid]
	if !exist {
		return nil, service.DbErrorNotFound(errors.New("the hosting doesn't exist"))
	}
	delete(h.store, uuid)
	delete(h.idxByName, hosting.Name)
	return hosting, nil
}
