package repository

import (
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/theskyinflames/cdmon2/app"
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

func populateConfig() *config.Config {
	return &config.Config{
		MinimalNumberOfCores:  1,
		MinimalSizeOfMemoryMb: 1,
		MinimalSizeOfDiskMb:   1,
	}
}

func TestHostingRepostitoryMap_Get(t *testing.T) {

	type fields struct {
		Mutex sync.Mutex
		store *StoreMock
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
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1}, nil
					},
				},
			},
			args:    args{uuid: "uuid1"},
			want:    &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
			wantErr: false,
		},
		{
			name: "given repository, when a not existing hosting is required, then it fails",
			fields: fields{
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return nil, app.DbErrorNotFound
					},
				},
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

	sliceOfhostings := []domain.Hosting{
		domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
		domain.Hosting{UUID: "uuid2", Name: "h2", Cores: 1, MemoryMb: 1, DiskMb: 1},
		domain.Hosting{UUID: "uuid3", Name: "h3", Cores: 1, MemoryMb: 1, DiskMb: 1},
	}

	type fields struct {
		store *StoreMock
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
				store: &StoreMock{
					GetAllFunc: func(pattern string, emptyRecordFunc config.EmptyRecordFunc) ([]interface{}, error) {
						return []interface{}{
							&domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
							&domain.Hosting{UUID: "uuid2", Name: "h2", Cores: 1, MemoryMb: 1, DiskMb: 1},
							&domain.Hosting{UUID: "uuid3", Name: "h3", Cores: 1, MemoryMb: 1, DiskMb: 1},
						}, nil
					},
				},
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

	type fields struct {
		Mutex sync.Mutex
		cfg   *config.Config
		store *StoreMock
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
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return nil, app.DbErrorNotFound
					},
					SetFunc: func(key string, item interface{}) error {
						return nil
					},
				},
			},
			args: args{
				hosting: &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
			},
			wantErr: false,
		},
		{
			name: "given a populated repository, when a hosting with a existing UUID is tried to be inserted, then if fails",
			fields: fields{
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return nil, app.DbErrorAlreadyExist
					},
					SetFunc: func(key string, item interface{}) error {
						return nil
					},
				},
			},
			args: args{
				hosting: &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
			},
			wantErr: true,
		},
		{
			name: "given a populated repository, when a hosting with a existing Name is tried to be inserted, then if fails",
			fields: fields{
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return nil, app.DbErrorAlreadyExist
					},
					SetFunc: func(key string, item interface{}) error {
						return nil
					},
				},
			},
			args: args{
				hosting: &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			h := HostingRepostitoryMap{
				Mutex: tt.fields.Mutex,
				cfg:   tt.fields.cfg,
				store: tt.fields.store,
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

	type fields struct {
		Mutex sync.Mutex
		cfg   *config.Config
		store *StoreMock
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
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 22, MemoryMb: 1, DiskMb: 1}, nil
					},
					SetFunc: func(key string, item interface{}) error {
						return nil
					},
				},
			},
			args: args{
				hosting: &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
			},
			wantErr: false,
		},
		{
			name: "given a populated repository, when a hosting with a not existing UUID is tried to be updated, then if fails",
			fields: fields{
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return nil, app.DbErrorNotFound
					},
					SetFunc: func(key string, item interface{}) error {
						return nil
					},
				},
			},
			args: args{
				hosting: &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
			},
			wantErr: true,
		},
		{
			name: "given a populated repository, when a hosting with a existing Name is tried to be updated, then if fails",
			fields: fields{
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						if key == "h2" {
							return "0", nil
						}
						return &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1}, nil
					},
					SetFunc: func(key string, item interface{}) error {
						return nil
					},
				},
			},
			args: args{
				hosting: &domain.Hosting{UUID: "uuid1", Name: "h2", Cores: 1, MemoryMb: 1, DiskMb: 1},
			},
			wantErr: true,
		},
	}
	for z, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			h := HostingRepostitoryMap{
				Mutex: tt.fields.Mutex,
				cfg:   tt.fields.cfg,
				store: tt.fields.store,
			}
			if err = h.Update(tt.args.hosting); (err != nil) != tt.wantErr {
				t.Errorf("HostingRepostitoryMap.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			switch z {
			case 0:
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

	type fields struct {
		Mutex sync.Mutex
		cfg   *config.Config
		store *StoreMock
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
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 22, MemoryMb: 1, DiskMb: 1}, nil
					},
					RemoveFunc: func(key string) error {
						return nil
					},
				},
			},
			args: args{
				uuid: domain.UUID("uuid1"),
			},
			want:    &domain.Hosting{UUID: "uuid1", Name: "h1", Cores: 22, MemoryMb: 1, DiskMb: 1},
			wantErr: false,
		},
		{
			name: "given a populated repository, when a hosting with not existing UUID is tried to be removed, then it fails",
			fields: fields{
				Mutex: sync.Mutex{},
				cfg:   cfg,
				store: &StoreMock{
					GetFunc: func(key string, item interface{}) (interface{}, error) {
						return nil, app.DbErrorNotFound
					},
					RemoveFunc: func(key string) error {
						return nil
					},
				},
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
				Mutex: tt.fields.Mutex,
				cfg:   tt.fields.cfg,
				store: tt.fields.store,
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
