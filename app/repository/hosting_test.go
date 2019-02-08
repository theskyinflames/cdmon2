package repository

import (
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	service "github.com/theskyinflames/cdmon2/app"
	"github.com/theskyinflames/cdmon2/app/config"
	"github.com/theskyinflames/cdmon2/app/domain"
)

type (
	ByUUID []domain.Hosting
)

func (a ByUUID) Len() int           { return len(a) }
func (a ByUUID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUUID) Less(i, j int) bool { return a[i].UUID < a[j].UUID }

func populateFixturesByUUID() map[domain.UUID]*domain.Hosting {
	return map[domain.UUID]*domain.Hosting{
		domain.UUID("uuid1"): &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
		domain.UUID("uuid2"): &domain.Hosting{UUID: "uuid2", Name: "h2", Cores: 1, MemoryMb: 1, DiskMb: 1},
		domain.UUID("uuid3"): &domain.Hosting{UUID: "uuid3", Name: "h3", Cores: 1, MemoryMb: 1, DiskMb: 1},
	}
}
func populateFixturesByName() map[string]*domain.Hosting {
	return map[string]*domain.Hosting{
		"h1": &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
		"h2": &domain.Hosting{UUID: "uuid2", Name: "h2", Cores: 1, MemoryMb: 1, DiskMb: 1},
		"h3": &domain.Hosting{UUID: "uuid3", Name: "h3", Cores: 1, MemoryMb: 1, DiskMb: 1},
	}
}

func populateConfig() *config.Config {
	return &config.Config{
		MinimalNumberOfCores:  1,
		MinimalSizeOfMemoryMb: 1,
		MinimalSizeOfDiskMb:   1,
	}
}

func TestHostingRepostitoryMap_Get(t *testing.T) {

	fixtures := populateFixturesByUUID()

	type fields struct {
		Mutex sync.Mutex
		store map[domain.UUID]*domain.Hosting
	}
	type args struct {
		uuid domain.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Hosting
		wantErr bool
	}{
		{
			name: "given repository, when an existing hosting is required, then it's returned",
			fields: fields{
				store: fixtures,
			},
			args:    args{uuid: "uuid1"},
			want:    fixtures[domain.UUID("uuid1")],
			wantErr: false,
		},
		{
			name: "given repository, when a not existing hosting is required, then it fails",
			fields: fields{
				store: fixtures,
			},
			args:    args{uuid: "uuid111"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HostingRepostitoryMap{
				store: tt.fields.store,
			}
			got, err := h.Get(tt.args.uuid)

			if (err != nil) != tt.wantErr {
				t.Errorf("HostingRepostitoryMap.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HostingRepostitoryMap.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostingRepostitoryMap_GetAll(t *testing.T) {

	fixtures := populateFixturesByUUID()

	sliceOfhostings := make([]domain.Hosting, len(fixtures))
	z := 0
	for _, v := range fixtures {
		sliceOfhostings[z] = *v
		z++
	}

	type fields struct {
		store map[domain.UUID]*domain.Hosting
	}
	tests := []struct {
		name    string
		fields  fields
		want    []domain.Hosting
		wantErr bool
	}{
		{
			name: "given a repository, when the list of hostins is required, then it's returned",
			fields: fields{
				store: fixtures,
			},
			want:    sliceOfhostings,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HostingRepostitoryMap{
				store: tt.fields.store,
			}
			got, err := h.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("HostingRepostitoryMap.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Before comparte the two slices, it must be sorted
			assert.Equal(t, len(got), len(tt.want))
			sort.Sort(ByUUID(got))
			sort.Sort(ByUUID(tt.want))

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HostingRepostitoryMap.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostingRepostitoryMap_Insert(t *testing.T) {

	cfg := populateConfig()
	fixtures := populateFixturesByUUID()
	fixturesByName := populateFixturesByName()
	newHosting := &domain.Hosting{UUID: "uuid33", Name: "h33", Cores: 1, MemoryMb: 1, DiskMb: 1}
	alreadyExistingUUID := &domain.Hosting{UUID: "uuid1", Name: "h134", Cores: 1, MemoryMb: 1, DiskMb: 1}
	alreadyExistingName := &domain.Hosting{UUID: "uuid44", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1}

	type fields struct {
		Mutex     sync.Mutex
		cfg       *config.Config
		store     map[domain.UUID]*domain.Hosting
		idxByName map[string]*domain.Hosting
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
			name: "given a populated repository, when a new hosting is inserted then all works fine",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				hosting: newHosting,
			},
			wantErr: false,
		},
		{
			name: "given a populated repository, when a hosting with a existing UUID is tried to be inserted, then if fails",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				hosting: alreadyExistingUUID,
			},
			wantErr: true,
		},
		{
			name: "given a populated repository, when a hosting with a existing Name is tried to be inserted, then if fails",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				hosting: alreadyExistingName,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			h := HostingRepostitoryMap{
				Mutex:     tt.fields.Mutex,
				cfg:       tt.fields.cfg,
				store:     tt.fields.store,
				idxByName: tt.fields.idxByName,
			}
			if err = h.Insert(tt.args.hosting); (err != nil) != tt.wantErr {
				t.Errorf("HostingRepostitoryMap.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				switch errors.Cause(err) {
				case service.DbErrorNotFound:
					t.Fail()
				case service.DbErrorAlreadyExist:
					t.Logf("OK")
				default:
					t.Fail()
				}
			}
		})
	}
}

func TestHostingRepostitoryMap_Update(t *testing.T) {
	cfg := populateConfig()
	fixtures := populateFixturesByUUID()
	fixturesByName := populateFixturesByName()
	updatedHosting := &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 22, MemoryMb: 1, DiskMb: 1}
	notExistingUUID := &domain.Hosting{UUID: "uuid133", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1}
	alreadyExistingName := &domain.Hosting{UUID: "uuid1", Name: "h2", Cores: 1, MemoryMb: 1, DiskMb: 1}

	type fields struct {
		Mutex     sync.Mutex
		cfg       *config.Config
		store     map[domain.UUID]*domain.Hosting
		idxByName map[string]*domain.Hosting
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
			name: "given a populated repository, when a hosting is updated then all works fine",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				hosting: updatedHosting,
			},
			wantErr: false,
		},
		{
			name: "given a populated repository, when a hosting with a not existing UUID is tried to be updated, then if fails",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				hosting: notExistingUUID,
			},
			wantErr: true,
		},
		{
			name: "given a populated repository, when a hosting with a existing Name is tried to be updated, then if fails",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				hosting: alreadyExistingName,
			},
			wantErr: true,
		},
	}
	for z, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			h := HostingRepostitoryMap{
				Mutex:     tt.fields.Mutex,
				cfg:       tt.fields.cfg,
				store:     tt.fields.store,
				idxByName: tt.fields.idxByName,
			}
			if err = h.Update(tt.args.hosting); (err != nil) != tt.wantErr {
				t.Errorf("HostingRepostitoryMap.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			switch z {
			case 0:
				assert.Equal(t, 22, h.store[domain.UUID("uuid1")].Cores)
			case 1:
				switch errors.Cause(err) {
				case service.DbErrorNotFound:
					t.Log(err.Error())
				default:
					t.Fail()
				}
			case 2:
				switch errors.Cause(err) {
				case service.DbErrorAlreadyExist:
					t.Log(err.Error())
				default:
					t.Fail()
				}
			default:
				t.Fail()
			}
		})
	}
}

func TestHostingRepostitoryMap_Remove(t *testing.T) {

	cfg := populateConfig()
	fixtures := populateFixturesByUUID()
	fixturesByName := populateFixturesByName()
	removedHosting := &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1}
	removedHostingUUID := domain.UUID("uuid1")
	afterRemoveFixtures := populateFixturesByUUID()
	delete(afterRemoveFixtures, domain.UUID("uuid1"))
	afterRemoveFixturesByName := populateFixturesByName()
	delete(afterRemoveFixturesByName, "h1")

	type fields struct {
		Mutex     sync.Mutex
		cfg       *config.Config
		store     map[domain.UUID]*domain.Hosting
		idxByName map[string]*domain.Hosting
	}
	type args struct {
		uuid domain.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Hosting
		wantErr bool
	}{
		{
			name: "given a populated repository, when a hosting is removed then all works fine",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				uuid: removedHostingUUID,
			},
			want:    removedHosting,
			wantErr: false,
		},
		{
			name: "given a populated repository, when a hosting with not existing UUID is tried to be removed, then it fails",
			fields: fields{
				Mutex:     sync.Mutex{},
				cfg:       cfg,
				store:     fixtures,
				idxByName: fixturesByName,
			},
			args: args{
				uuid: domain.UUID("uuid9999"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for z, tt := range tests {
		var (
			err error
			got *domain.Hosting
		)
		t.Run(tt.name, func(t *testing.T) {
			h := HostingRepostitoryMap{
				Mutex:     tt.fields.Mutex,
				cfg:       tt.fields.cfg,
				store:     tt.fields.store,
				idxByName: tt.fields.idxByName,
			}
			got, err = h.Remove(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("HostingRepostitoryMap.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HostingRepostitoryMap.Remove() = %v, want %v", got, tt.want)
			}

			switch z {
			case 0:
				if !reflect.DeepEqual(h.store, afterRemoveFixtures) {
					t.Errorf("HostingRepostitoryMap.Remove() = %v, want %v", h.store, afterRemoveFixtures)
				}
				if !reflect.DeepEqual(h.idxByName, afterRemoveFixturesByName) {
					t.Errorf("HostingRepostitoryMap.Remove() = %v, want %v", h.idxByName, afterRemoveFixturesByName)
				}
			case 1:
				switch errors.Cause(err) {
				case service.DbErrorAlreadyExist:
					t.Fail()
				case service.DbErrorNotFound:
					t.Logf(err.Error())
				default:
					t.Fail()
				}
			}
		})
	}
}
