package domain

import (
	"github.com/pkg/errors"
	gouuid "github.com/satori/go.uuid"
	"github.com/theskyinflames/cdmon2/hosting/config"
)

const ()

type (
	Hosting struct {
		UUID     UUID   `json:"uuid"`
		Name     string `json:"name"`
		Cores    int    `json:"cores"`
		MemoryMb int    `json:"memorymb"`
		DiskMb   int    `json:"diskmb"`
	}
)

func (u UUID) Validate() error {
	if len(string(u)) == 0 {
		return errors.New("UUID can't be empty")
	}
	return nil
}

func NewHosting(name string, cores, memorymb, diskmb int) (*Hosting, error) {

	uuid := gouuid.NewV1()
	return &Hosting{
		UUID:     UUID(uuid.String()),
		Name:     name,
		Cores:    cores,
		MemoryMb: memorymb,
		DiskMb:   diskmb,
	}, nil
}

func (h *Hosting) Validate(cfg *config.Config) (err error) {

	err = h.UUID.Validate()
	if err != nil {
		return err
	}

	if len(h.Name) == 0 {
		return errors.New("Name can't be empty")
	}

	if h.Cores < cfg.MinimalNumberOfCores {
		return errors.New("Number of cores can't be zero")
	}

	if h.MemoryMb < cfg.MinimalSizeOfMemoryMb {
		return errors.New("Memory size can't be less than 1 Mb")
	}

	if h.DiskMb < cfg.MinimalSizeOfDiskMb {
		return errors.New("Disk space size cant't be less than 1 Mb")
	}

	return nil
}
