package service

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/theskyinflames/cdmon2/app/config"
	"github.com/theskyinflames/cdmon2/app/domain"
)

func populateConfig() *config.Config {
	return &config.Config{
		MinimalNumberOfCores:  1,
		MinimalSizeOfMemoryMb: 1,
		MinimalSizeOfDiskMb:   1,
	}
}

func populateHostings() []domain.Hosting {
	return []domain.Hosting{
		domain.Hosting{UUID: domain.UUID("uuid1"), Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
		domain.Hosting{UUID: domain.UUID("uuid2"), Name: "h2", Cores: 1, MemoryMb: 1, DiskMb: 1},
		domain.Hosting{UUID: domain.UUID("uuid3"), Name: "h3", Cores: 1, MemoryMb: 1, DiskMb: 1},
	}
}

func NewHostingRepositoryMockOK() *HostingRepositoryMock {
	return &HostingRepositoryMock{
		InsertFunc: func(hosting *domain.Hosting) error {
			return nil
		},
		RemoveFunc: func(uuid domain.UUID) (*domain.Hosting, error) {
			return &populateHostings()[0], nil
		},
		GetAllFunc: func() ([]domain.Hosting, error) {
			return populateHostings(), nil
		},
		GetFunc: func(uuid domain.UUID) (*domain.Hosting, error) {
			return &populateHostings()[0], nil
		},
		UpdateFunc: func(hosting *domain.Hosting) error {
			return nil
		},
	}
}
func NewServerDomainMockOK() *ServerDomainMock {
	return &ServerDomainMock{
		AddHostingFunc: func(hosting *domain.Hosting, cfg *config.Config) error {
			return nil
		},
		RemoveHostingFunc: func(hosting *domain.Hosting) error {
			return nil
		},
		UpdateHostingFunc: func(hosting, old *domain.Hosting, cfg *config.Config) error {
			return nil
		},
	}
}

func TestServerService_CreateHosting(t *testing.T) {

	cfg := populateConfig()
	log := logrus.New()

	type fields struct {
		Mutex             sync.Mutex
		log               *logrus.Logger
		cfg               *config.Config
		hostingRepository *HostingRepositoryMock
		serverDomain      *ServerDomainMock
	}
	type args struct {
		name     string
		cores    int
		memorymb int
		diskmb   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "given a server, when a new hosting is created, then all works fine",
			fields: fields{
				Mutex:             sync.Mutex{},
				log:               log,
				cfg:               cfg,
				hostingRepository: NewHostingRepositoryMockOK(),
				serverDomain:      NewServerDomainMockOK(),
			},
			args:    args{name: "h1", cores: 1, memorymb: 1, diskmb: 1},
			wantErr: false,
		},
		{
			name: "given a server, when a new hosting is created and the repository fails, then it fails",
			fields: fields{
				Mutex: sync.Mutex{},
				log:   log,
				cfg:   cfg,
				hostingRepository: &HostingRepositoryMock{
					InsertFunc: func(hosting *domain.Hosting) error {
						return errors.New("random error")
					},
				},
				serverDomain: NewServerDomainMockOK(),
			},
			args:    args{name: "h1", cores: 1, memorymb: 1, diskmb: 1},
			wantErr: true,
		},
		{
			name: "given a server, when a new hosting is created and the server domains call fails, then it fails",
			fields: fields{
				Mutex:             sync.Mutex{},
				log:               log,
				cfg:               cfg,
				hostingRepository: NewHostingRepositoryMockOK(),
				serverDomain: &ServerDomainMock{
					AddHostingFunc: func(hosting *domain.Hosting, cfg *config.Config) error {
						return errors.New("random error")
					},
					RemoveHostingFunc: func(hosting *domain.Hosting) error {
						return nil
					},
				},
			},
			args:    args{name: "h1", cores: 1, memorymb: 1, diskmb: 1},
			wantErr: true,
		},
	}
	for z, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got domain.UUID
				err error
			)
			s := &ServerService{
				Mutex:             tt.fields.Mutex,
				log:               tt.fields.log,
				cfg:               tt.fields.cfg,
				hostingRepository: tt.fields.hostingRepository,
				serverDomain:      tt.fields.serverDomain,
			}
			got, err = s.CreateHosting(tt.args.name, tt.args.cores, tt.args.memorymb, tt.args.diskmb)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerService.CreateHosting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				assert.NoError(t, got.Validate())
			}

			// Validate the execution flux for each scenery
			switch z {
			case 0:
				assert.Equal(t, 1, len(tt.fields.serverDomain.AddHostingCalls()))
				assert.Equal(t, 1, len(tt.fields.hostingRepository.InsertCalls()))
				assert.Equal(t, 0, len(tt.fields.serverDomain.RemoveHostingCalls()))
			case 1:
				assert.Equal(t, 1, len(tt.fields.serverDomain.AddHostingCalls()))
				assert.Equal(t, 1, len(tt.fields.hostingRepository.InsertCalls()))
				assert.Equal(t, 1, len(tt.fields.serverDomain.RemoveHostingCalls()))
			case 2:
				assert.Equal(t, 1, len(tt.fields.serverDomain.AddHostingCalls()))
				assert.Equal(t, 0, len(tt.fields.hostingRepository.InsertCalls()))
				assert.Equal(t, 0, len(tt.fields.serverDomain.RemoveHostingCalls()))
			}
		})
	}
}

func TestServerService_GetHostings(t *testing.T) {

	cfg := populateConfig()
	hostings := populateHostings()
	log := logrus.New()

	type fields struct {
		Mutex             sync.Mutex
		log               *logrus.Logger
		cfg               *config.Config
		hostingRepository *HostingRepositoryMock
		serverDomain      *ServerDomainMock
	}
	tests := []struct {
		name    string
		fields  fields
		want    []domain.Hosting
		wantErr bool
	}{
		{
			name: "given a server, when the list of hosting is requested, then it's returned",
			fields: fields{
				Mutex:             sync.Mutex{},
				log:               log,
				cfg:               cfg,
				hostingRepository: NewHostingRepositoryMockOK(),
				serverDomain:      nil,
			},
			want:    hostings,
			wantErr: false,
		},
		{
			name: "given a server, when the list of hosting is requested and the repository fails, then it fails",
			fields: fields{
				Mutex: sync.Mutex{},
				log:   log,
				cfg:   cfg,
				hostingRepository: &HostingRepositoryMock{
					GetAllFunc: func() ([]domain.Hosting, error) {
						return nil, errors.New("random error")
					},
				},
				serverDomain: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerService{
				Mutex:             tt.fields.Mutex,
				log:               tt.fields.log,
				cfg:               tt.fields.cfg,
				hostingRepository: tt.fields.hostingRepository,
				serverDomain:      tt.fields.serverDomain,
			}
			got, err := s.GetHostings()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerService.GetHostings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerService.GetHostings() = %v, want %v", got, tt.want)
			}

			assert.Equal(t, 1, len(tt.fields.hostingRepository.GetAllCalls()))
		})
	}
}

func TestServerService_RemoveHosting(t *testing.T) {

	cfg := populateConfig()
	hosting := populateHostings()[0]
	log := logrus.New()

	type fields struct {
		Mutex             sync.Mutex
		log               *logrus.Logger
		cfg               *config.Config
		hostingRepository *HostingRepositoryMock
		serverDomain      *ServerDomainMock
	}
	type args struct {
		uuid domain.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "given a server, when a hosting is removed, then all works fine",
			fields: fields{
				Mutex:             sync.Mutex{},
				log:               log,
				cfg:               cfg,
				hostingRepository: NewHostingRepositoryMockOK(),
				serverDomain:      NewServerDomainMockOK(),
			},
			args:    args{uuid: hosting.UUID},
			wantErr: false,
		},
		{
			name: "given a server, when a hosting is going to be removed and the repository fails, then it fails",
			fields: fields{
				Mutex: sync.Mutex{},
				log:   log,
				cfg:   cfg,
				hostingRepository: &HostingRepositoryMock{
					InsertFunc: func(hosting *domain.Hosting) error {
						return nil
					},
					RemoveFunc: func(uuid domain.UUID) (*domain.Hosting, error) {
						return nil, errors.New("random error")
					},
				},
				serverDomain: NewServerDomainMockOK(),
			},
			args:    args{uuid: hosting.UUID},
			wantErr: true,
		},
		{
			name: "given a server, when a hosting is going to be removed and the server domains fails, then it fails",
			fields: fields{
				Mutex:             sync.Mutex{},
				log:               log,
				cfg:               cfg,
				hostingRepository: NewHostingRepositoryMockOK(),
				serverDomain: &ServerDomainMock{
					RemoveHostingFunc: func(hosting *domain.Hosting) error {
						return errors.New("random error")
					},
				},
			},
			args:    args{uuid: hosting.UUID},
			wantErr: true,
		},
	}
	for z, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerService{
				Mutex:             tt.fields.Mutex,
				log:               tt.fields.log,
				cfg:               tt.fields.cfg,
				hostingRepository: tt.fields.hostingRepository,
				serverDomain:      tt.fields.serverDomain,
			}
			if err := s.RemoveHosting(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ServerService.RemoveHosting() error = %v, wantErr %v", err, tt.wantErr)
			}

			switch z {
			case 0:
				assert.Equal(t, 1, len(tt.fields.hostingRepository.RemoveCalls()))
				assert.Equal(t, 0, len(tt.fields.hostingRepository.InsertCalls()))
				assert.Equal(t, 1, len(tt.fields.serverDomain.RemoveHostingCalls()))
			case 1:
				assert.Equal(t, 1, len(tt.fields.hostingRepository.RemoveCalls()))
				assert.Equal(t, 0, len(tt.fields.hostingRepository.InsertCalls()))
				assert.Equal(t, 0, len(tt.fields.serverDomain.RemoveHostingCalls()))
			case 2:
				assert.Equal(t, 1, len(tt.fields.hostingRepository.RemoveCalls()))
				assert.Equal(t, 1, len(tt.fields.hostingRepository.InsertCalls()))
				assert.Equal(t, 1, len(tt.fields.serverDomain.RemoveHostingCalls()))
			}
		})
	}
}

func TestServerService_UpdateHosting(t *testing.T) {

	cfg := populateConfig()
	hosting := populateHostings()[0]
	log := logrus.New()

	type fields struct {
		Mutex             sync.Mutex
		log               *logrus.Logger
		cfg               *config.Config
		hostingRepository *HostingRepositoryMock
		serverDomain      *ServerDomainMock
	}
	type args struct {
		hosting *domain.Hosting
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "given a server, when a hosting is updated, then all workds fine",
			fields: fields{
				Mutex:             sync.Mutex{},
				log:               log,
				cfg:               cfg,
				hostingRepository: NewHostingRepositoryMockOK(),
				serverDomain:      NewServerDomainMockOK(),
			},
			args:    args{hosting: &hosting},
			wantErr: false,
		},
		{
			name: "given a server, when a hosting is tried to be updated and the repository.Get op fails, then it fails",
			fields: fields{
				Mutex: sync.Mutex{},
				log:   log,
				cfg:   cfg,
				hostingRepository: &HostingRepositoryMock{
					GetFunc: func(uuid domain.UUID) (*domain.Hosting, error) {
						return nil, errors.New("random error")
					},
				},
				serverDomain: NewServerDomainMockOK(),
			},
			args:    args{hosting: &hosting},
			wantErr: true,
		},
		{
			name: "given a server, when a existing hosting is tried to be updated and the server domain fails, then it fails",
			fields: fields{
				Mutex:             sync.Mutex{},
				log:               log,
				cfg:               cfg,
				hostingRepository: NewHostingRepositoryMockOK(),
				serverDomain: &ServerDomainMock{
					UpdateHostingFunc: func(hosting, old *domain.Hosting, cfg *config.Config) error {
						return errors.New("random error")
					},
				},
			},
			args:    args{hosting: &hosting},
			wantErr: true,
		},
		{
			name: "given a server, when a existing hosting is tried to be updated and the repository.Update op fails, then it fails",
			fields: fields{
				Mutex: sync.Mutex{},
				log:   log,
				cfg:   cfg,
				hostingRepository: &HostingRepositoryMock{
					GetFunc: func(uuid domain.UUID) (*domain.Hosting, error) {
						return &hosting, nil
					},
					UpdateFunc: func(hosting *domain.Hosting) error {
						return errors.New("random error")
					},
				},
				serverDomain: NewServerDomainMockOK(),
			},
			args:    args{hosting: &hosting},
			wantErr: true,
		},
	}
	for z, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerService{
				Mutex:             tt.fields.Mutex,
				log:               tt.fields.log,
				cfg:               tt.fields.cfg,
				hostingRepository: tt.fields.hostingRepository,
				serverDomain:      tt.fields.serverDomain,
			}
			if err := s.UpdateHosting(tt.args.hosting); (err != nil) != tt.wantErr {
				t.Errorf("ServerService.UpdateHosting() error = %v, wantErr %v", err, tt.wantErr)
			}

			switch z {
			case 0:
				assert.Equal(t, 1, len(tt.fields.hostingRepository.GetCalls()))
				assert.Equal(t, 1, len(tt.fields.serverDomain.UpdateHostingCalls()))
				assert.Equal(t, 1, len(tt.fields.hostingRepository.UpdateCalls()))
			case 1:
				assert.Equal(t, 1, len(tt.fields.hostingRepository.GetCalls()))
				assert.Equal(t, 0, len(tt.fields.serverDomain.UpdateHostingCalls()))
				assert.Equal(t, 0, len(tt.fields.hostingRepository.UpdateCalls()))
			case 2:
				assert.Equal(t, 1, len(tt.fields.hostingRepository.GetCalls()))
				assert.Equal(t, 1, len(tt.fields.serverDomain.UpdateHostingCalls()))
				assert.Equal(t, 0, len(tt.fields.hostingRepository.UpdateCalls()))
			case 3:
				assert.Equal(t, 1, len(tt.fields.hostingRepository.GetCalls()))
				assert.Equal(t, 1, len(tt.fields.serverDomain.UpdateHostingCalls()))
				assert.Equal(t, 1, len(tt.fields.hostingRepository.UpdateCalls()))
			}
		})
	}
}
