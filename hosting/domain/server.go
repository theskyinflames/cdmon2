package domain

import (
	"sync"

	"github.com/pkg/errors"
	gouuid "github.com/satori/go.uuid"
	"github.com/theskyinflames/cdmon2/hosting/config"
)

type (
	UUID string

	Server struct {
		sync.Mutex
		UUID                    UUID `json:"uuid"`
		TotalCores              int  `json:"total_cores"`
		TotalSizeOfMemoryMb     int  `json:"total_memory_mb"`
		TotalSizeOfDiskMb       int  `json:"total_disk_mb"`
		AvailableCores          int  `json:"available_cores"`
		AvailableSizeOfMemoryMb int  `json:"available_memory_mb"`
		AvailableSizeOfDiskMb   int  `json:"available_disk_mb`
	}
)

func NewServer(cfg *config.Config) (*Server, error) {

	uuid := gouuid.NewV1()
	return &Server{
		UUID:                    UUID(uuid.String()),
		TotalCores:              cfg.TotalNumberOfCores,
		TotalSizeOfMemoryMb:     cfg.TotalSizeOfMemoryMb,
		TotalSizeOfDiskMb:       cfg.TotalSizeOfDiskMb,
		AvailableCores:          cfg.TotalNumberOfCores,
		AvailableSizeOfMemoryMb: cfg.TotalSizeOfMemoryMb,
		AvailableSizeOfDiskMb:   cfg.TotalSizeOfDiskMb,
	}, nil
}

func (s *Server) CreateHosting(hosting *Hosting, cfg *config.Config) error {
	err := s.checkForResourcesAvailability(hosting)
	if err != nil {
		return errors.Wrap(err, "there aren't resources enough to create the hosting")
	}

	s.assignResources(hosting)
	return nil
}

func (s *Server) checkForResourcesAvailability(hosting *Hosting) error {
	if (s.AvailableCores - hosting.Cores) < 0 {
		return errors.New("there is not cores enough")
	}
	if (s.AvailableSizeOfMemoryMb - hosting.MemoryMb) < 0 {
		return errors.New("there is not memory mb enough")
	}
	if (s.AvailableSizeOfDiskMb - hosting.DiskMb) < 0 {
		return errors.New("there is not disk space mb enough")
	}
	return nil
}

func (s *Server) assignResources(hosting *Hosting) {
	s.AvailableCores -= hosting.Cores
	s.AvailableSizeOfMemoryMb -= hosting.MemoryMb
	s.AvailableSizeOfDiskMb -= hosting.DiskMb
	return
}

func (s *Server) RemoveHosting(hosting *Hosting) error {
	s.restoreResources(hosting)
	return nil
}

func (s *Server) restoreResources(hosting *Hosting) {
	s.AvailableCores += hosting.Cores
	s.AvailableSizeOfMemoryMb += hosting.MemoryMb
	s.AvailableSizeOfDiskMb += hosting.DiskMb
}

func (s *Server) UpdateHosting(hosting, old *Hosting, cfg *config.Config) error {
	s.restoreResources(old)

	err := s.checkForResourcesAvailability(hosting)
	if err != nil {
		s.assignResources(old)
		return errors.Wrapf(err, "somthing went wrong when trying to update the hosting %s", string(hosting.UUID))
	}

	s.assignResources(hosting)
	return nil
}
