package repository

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/theskyinflames/cdmon2/app"
	"github.com/theskyinflames/cdmon2/app/config"
	"github.com/theskyinflames/cdmon2/app/domain"
)

type (
	Store interface {
		Connect() error
		Close() error
		Get(key string, item interface{}) (interface{}, error)
		GetAll(pattern string, emptyRecordFunc config.EmptyRecordFunc) ([]interface{}, error)
		Set(key string, item interface{}) error
		Remove(key string) error
	}

	HostingRepostitoryMap struct {
		sync.Mutex
		cfg   *config.Config
		store Store
	}
)

func NewHostingReposytoryMap(cfg *config.Config, store Store) *HostingRepostitoryMap {
	return &HostingRepostitoryMap{
		Mutex: sync.Mutex{},
		cfg:   cfg,
		store: store,
	}
}

func (h HostingRepostitoryMap) Get(uuid domain.UUID) (*domain.Hosting, error) {
	var (
		err  error
		item interface{}
	)

	item, err = h.store.Get(string(uuid), &domain.Hosting{})
	if err != nil {
		switch err {
		case app.DbErrorNotFound:
			return nil, errors.Wrapf(err, "uuid: %s", string(uuid))
		default:
			return nil, err
		}
	}
	return item.(*domain.Hosting), nil
}

func (h HostingRepostitoryMap) GetAll() ([]domain.Hosting, error) {

	var emptyRecordFunc config.EmptyRecordFunc = func() interface{} {
		return &domain.Hosting{}
	}
	slice, err := h.store.GetAll("*-*", emptyRecordFunc)
	if err != nil {
		return nil, err
	}
	hostings := make([]domain.Hosting, len(slice))
	for z, v := range slice {
		hostings[z] = *v.(*domain.Hosting)
	}
	return hostings, nil
}

func (h HostingRepostitoryMap) Insert(hosting *domain.Hosting) error {
	_, err := h.store.Get(string(hosting.UUID), hosting)
	switch err {
	case nil:
		return errors.Wrapf(app.DbErrorAlreadyExist, "uuid: %s", string(hosting.UUID))
	case app.DbErrorNotFound:
	default:
		return err
	}

	var s string
	_, err = h.store.Get(hosting.Name, &s)
	switch err {
	case nil:
		return errors.Wrapf(app.DbErrorAlreadyExist, "name: %s", hosting.Name)
	case app.DbErrorNotFound:
	default:
		return err
	}

	err = hosting.Validate(h.cfg)
	if err != nil {
		return err
	}

	h.store.Set(string(hosting.UUID), *hosting)
	h.store.Set(hosting.Name, "0")
	return nil
}

func (h HostingRepostitoryMap) Update(hosting *domain.Hosting) error {
	old, err := h.Get(hosting.UUID)
	if err != nil {
		return err
	}
	if old.Name != hosting.Name {
		// Check for a already existing name
		var s string
		_, err = h.store.Get(hosting.Name, &s)
		switch err {
		case nil:
			return errors.Wrapf(app.DbErrorAlreadyExist, "name: %s", hosting.Name)
		case app.DbErrorNotFound:
		default:
			return err
		}
	}
	err = hosting.Validate(h.cfg)
	if err != nil {
		return err
	}

	h.store.Set(string(hosting.UUID), *hosting)
	h.store.Set(hosting.Name, "0")
	return nil
}

func (h HostingRepostitoryMap) Remove(uuid domain.UUID) (*domain.Hosting, error) {
	hosting, err := h.Get(uuid)
	if err != nil {
		return nil, err
	}
	h.store.Remove(string(uuid))
	h.store.Remove(hosting.Name)
	return hosting, nil
}
