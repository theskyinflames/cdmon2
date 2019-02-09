package service

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/theskyinflames/cdmon2/app/config"
	"github.com/theskyinflames/cdmon2/app/domain"
)

type (
	HostingRepository interface {
		Get(uuid domain.UUID) (*domain.Hosting, error)
		GetAll() ([]domain.Hosting, error)
		Insert(hosting *domain.Hosting) error
		Update(hosting *domain.Hosting) error
		Remove(uuid domain.UUID) (*domain.Hosting, error)
	}

	ServerDomain interface {
		AddHosting(hosting *domain.Hosting, cfg *config.Config) error
		UpdateHosting(hosting, old *domain.Hosting, cfg *config.Config) error
		RemoveHosting(hosting *domain.Hosting) error
	}

	ServerService struct {
		sync.Mutex
		log               *logrus.Logger
		cfg               *config.Config
		hostingRepository HostingRepository
		serverDomain      ServerDomain
	}
)

func NewServer(hostingRepository HostingRepository, serverDomain ServerDomain, cfg *config.Config, log *logrus.Logger) *ServerService {
	return &ServerService{
		Mutex:             sync.Mutex{},
		log:               log,
		cfg:               cfg,
		hostingRepository: hostingRepository,
		serverDomain:      serverDomain,
	}
}

func (s *ServerService) CreateHosting(name string, cores int, memorymb int, diskmb int) (domain.UUID, error) {
	hosting, err := domain.NewHosting(name, cores, memorymb, diskmb)
	if err != nil {
		return domain.UUID(""), err
	}

	s.Lock()
	defer s.Unlock()

	err = s.serverDomain.AddHosting(hosting, s.cfg)
	if err != nil {
		return domain.UUID(""), err
	}

	err = s.hostingRepository.Insert(hosting)
	if err != nil {
		s.serverDomain.RemoveHosting(hosting)
		return domain.UUID(""), err
	}

	s.log.WithFields(logrus.Fields{"uuid": string(hosting.UUID)}).Info("created hosting")
	return hosting.UUID, nil
}

func (s *ServerService) GetHostings() ([]domain.Hosting, error) {
	return s.hostingRepository.GetAll()
}

func (s *ServerService) RemoveHosting(uuid domain.UUID) error {
	s.Lock()
	defer s.Unlock()

	hosting, err := s.hostingRepository.Remove(uuid)
	if err != nil {
		return err
	}

	err = s.serverDomain.RemoveHosting(hosting)
	if err != nil {
		s.hostingRepository.Insert(hosting)
		return err
	}

	s.log.WithFields(logrus.Fields{"uuid": string(hosting.UUID)}).Info("removed hosting")
	return nil
}

func (s *ServerService) UpdateHosting(hosting *domain.Hosting) error {
	s.Lock()
	defer s.Unlock()

	old, err := s.hostingRepository.Get(hosting.UUID)
	if err != nil {
		return err
	}

	err = s.serverDomain.UpdateHosting(hosting, old, s.cfg)
	if err != nil {
		return err
	}

	err = s.hostingRepository.Update(hosting)
	if err != nil {
		return err
	}

	s.log.WithFields(logrus.Fields{"uuid": string(hosting.UUID)}).Info("updated hosting")
	return nil
}

func (s *ServerService) GetServerStatus() *domain.Server {
	return s.serverDomain.(*domain.Server)
}
