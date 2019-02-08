package domain

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theskyinflames/cdmon2/app/config"
)

func populateServer() *Server {
	return &Server{
		Mutex:                   sync.Mutex{},
		UUID:                    UUID("uuid1"),
		TotalCores:              100,
		TotalSizeOfMemoryMb:     100,
		TotalSizeOfDiskMb:       100,
		AvailableCores:          100,
		AvailableSizeOfMemoryMb: 100,
		AvailableSizeOfDiskMb:   100,
	}
}

func populateHosting(cores, memorymb, diskmb int) *Hosting {
	return &Hosting{UUID: UUID("uuid"), Name: "h1", Cores: cores, MemoryMb: memorymb, DiskMb: diskmb}
}

func TestServer_AddHosting(t *testing.T) {

	server := populateServer()
	hosting := populateHosting(1, 1, 1)
	overSizedCoresHosting := populateHosting(101, 1, 1)
	overSizedMemoryHosting := populateHosting(1, 101, 1)
	overSizedDiskHosting := populateHosting(1, 1, 101)

	cfg := populateConfig()

	type fields struct {
		server *Server
	}
	type args struct {
		hosting *Hosting
		cfg     *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "given a server, when a hosting is added, then works fine",
			fields: fields{
				server: server,
			},
			args: args{
				hosting: hosting,
				cfg:     cfg,
			},
			wantErr: false,
		},
		{
			name: "given a server, when a over sized in cores hosting is added, then it fails",
			fields: fields{
				server: server,
			},
			args: args{
				hosting: overSizedCoresHosting,
				cfg:     cfg,
			},
			wantErr: true,
		},
		{
			name: "given a server, when a over sized in memory hosting is added, then it fails",
			fields: fields{
				server: server,
			},
			args: args{
				hosting: overSizedMemoryHosting,
				cfg:     cfg,
			},
			wantErr: true,
		},
		{
			name: "given a server, when a over sized in disk hosting is added, then it fails",
			fields: fields{
				server: server,
			},
			args: args{
				hosting: overSizedDiskHosting,
				cfg:     cfg,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.server
			if err := s.AddHosting(tt.args.hosting, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Server.AddHosting() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, 99, s.AvailableCores)
			assert.Equal(t, 99, s.AvailableSizeOfMemoryMb)
			assert.Equal(t, 99, s.AvailableSizeOfDiskMb)
		})
	}
}

func TestServer_RemoveHosting(t *testing.T) {

	hosting := populateHosting(1, 1, 1)
	notExistHosting := populateHosting(1, 1, 1)
	notExistHosting.UUID = UUID("uuid89")

	type fields struct {
		Mutex                   sync.Mutex
		UUID                    UUID
		TotalCores              int
		TotalSizeOfMemoryMb     int
		TotalSizeOfDiskMb       int
		AvailableCores          int
		AvailableSizeOfMemoryMb int
		AvailableSizeOfDiskMb   int
	}
	type args struct {
		hosting *Hosting
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "given a server, when an existing hosting is removed, all works fine",
			fields: fields{
				Mutex:                   sync.Mutex{},
				UUID:                    UUID("uuid1"),
				TotalCores:              100,
				TotalSizeOfMemoryMb:     100,
				TotalSizeOfDiskMb:       100,
				AvailableCores:          99,
				AvailableSizeOfMemoryMb: 99,
				AvailableSizeOfDiskMb:   99,
			},
			args: args{
				hosting: hosting,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Mutex:                   tt.fields.Mutex,
				UUID:                    tt.fields.UUID,
				TotalCores:              tt.fields.TotalCores,
				TotalSizeOfMemoryMb:     tt.fields.TotalSizeOfMemoryMb,
				TotalSizeOfDiskMb:       tt.fields.TotalSizeOfDiskMb,
				AvailableCores:          tt.fields.AvailableCores,
				AvailableSizeOfMemoryMb: tt.fields.AvailableSizeOfMemoryMb,
				AvailableSizeOfDiskMb:   tt.fields.AvailableSizeOfDiskMb,
			}
			if err := s.RemoveHosting(tt.args.hosting); (err != nil) != tt.wantErr {
				t.Errorf("Server.RemoveHosting() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, 100, s.AvailableCores)
			assert.Equal(t, 100, s.AvailableSizeOfMemoryMb)
			assert.Equal(t, 100, s.AvailableSizeOfDiskMb)
		})
	}
}

func TestServer_UpdateHosting(t *testing.T) {

	cfg := populateConfig()
	hostingOld := populateHosting(1, 1, 1)
	hostingNewOK := populateHosting(10, 1, 1)
	hostingNewOverSized := populateHosting(101, 1, 1)

	type fields struct {
		Mutex                   sync.Mutex
		UUID                    UUID
		TotalCores              int
		TotalSizeOfMemoryMb     int
		TotalSizeOfDiskMb       int
		AvailableCores          int
		AvailableSizeOfMemoryMb int
		AvailableSizeOfDiskMb   int
	}
	type args struct {
		hosting *Hosting
		old     *Hosting
		cfg     *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "given a server, when a hosting is updated and the new configuration fits in the server, then all works fine",
			fields: fields{
				Mutex:                   sync.Mutex{},
				UUID:                    UUID("uuid1"),
				TotalCores:              100,
				TotalSizeOfMemoryMb:     100,
				TotalSizeOfDiskMb:       100,
				AvailableCores:          99,
				AvailableSizeOfMemoryMb: 99,
				AvailableSizeOfDiskMb:   99,
			},
			args: args{
				hosting: hostingNewOK,
				old:     hostingOld,
				cfg:     cfg,
			},
			wantErr: false,
		},
		{
			name: "given a server, when a hosting is updated and the new configuration doesn't fit in the server, then it fails",
			fields: fields{
				Mutex:                   sync.Mutex{},
				UUID:                    UUID("uuid1"),
				TotalCores:              100,
				TotalSizeOfMemoryMb:     100,
				TotalSizeOfDiskMb:       100,
				AvailableCores:          99,
				AvailableSizeOfMemoryMb: 99,
				AvailableSizeOfDiskMb:   99,
			},
			args: args{
				hosting: hostingNewOverSized,
				old:     hostingOld,
				cfg:     cfg,
			},
			wantErr: true,
		},
	}
	for z, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Mutex:                   tt.fields.Mutex,
				UUID:                    tt.fields.UUID,
				TotalCores:              tt.fields.TotalCores,
				TotalSizeOfMemoryMb:     tt.fields.TotalSizeOfMemoryMb,
				TotalSizeOfDiskMb:       tt.fields.TotalSizeOfDiskMb,
				AvailableCores:          tt.fields.AvailableCores,
				AvailableSizeOfMemoryMb: tt.fields.AvailableSizeOfMemoryMb,
				AvailableSizeOfDiskMb:   tt.fields.AvailableSizeOfDiskMb,
			}
			if err := s.UpdateHosting(tt.args.hosting, tt.args.old, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Server.UpdateHosting() error = %v, wantErr %v", err, tt.wantErr)
			}

			switch z {
			case 0:
				assert.Equal(t, 90, s.AvailableCores)
				assert.Equal(t, 99, s.AvailableSizeOfMemoryMb)
				assert.Equal(t, 99, s.AvailableSizeOfDiskMb)
			case 1:
				assert.Equal(t, 99, s.AvailableCores)
				assert.Equal(t, 99, s.AvailableSizeOfMemoryMb)
				assert.Equal(t, 99, s.AvailableSizeOfDiskMb)
			}
		})
	}
}
