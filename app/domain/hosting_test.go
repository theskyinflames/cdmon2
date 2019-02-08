package domain

import (
	"reflect"
	"testing"

	"github.com/theskyinflames/cdmon2/app/config"
)

func populateConfig() *config.Config {
	return &config.Config{
		MinimalNumberOfCores:  1,
		MinimalSizeOfMemoryMb: 1,
		MinimalSizeOfDiskMb:   1,
	}
}
func TestUUID_Validate(t *testing.T) {
	tests := []struct {
		name    string
		u       UUID
		wantErr bool
	}{
		{
			name:    "given a valid UUID, when it's validated, then all works fine",
			u:       UUID("uuid1"),
			wantErr: false,
		},
		{
			name:    "given a empty UUID, when it's validated, then it fails",
			u:       UUID(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UUID.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHosting_Validate(t *testing.T) {

	cfg := populateConfig()

	type fields struct {
		UUID     UUID
		Name     string
		Cores    int
		MemoryMb int
		DiskMb   int
	}
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "given a valid hosting, when it's validated then all works fine",
			fields:  fields{UUID: UUID("uuid"), Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 1},
			args:    args{cfg: cfg},
			wantErr: false,
		},
		{
			name:    "given a hosting without a name, when it's validated then it fails",
			fields:  fields{UUID: UUID("uuid"), Name: "", Cores: 1, MemoryMb: 1, DiskMb: 1},
			args:    args{cfg: cfg},
			wantErr: true,
		},
		{
			name:    "given a hosting with a invalid number of cores, when it's validated then it fails",
			fields:  fields{UUID: UUID("uuid"), Name: "h1", Cores: 0, MemoryMb: 1, DiskMb: 1},
			args:    args{cfg: cfg},
			wantErr: true,
		},
		{
			name:    "given a hosting with a invalid memory size, when it's validated then it fails",
			fields:  fields{UUID: UUID("uuid"), Name: "h1", Cores: 1, MemoryMb: 0, DiskMb: 1},
			args:    args{cfg: cfg},
			wantErr: true,
		},
		{
			name:    "given a hosting with a invalid disk size, when it's validated then it fails",
			fields:  fields{UUID: UUID("uuid"), Name: "h1", Cores: 1, MemoryMb: 1, DiskMb: 0},
			args:    args{cfg: cfg},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Hosting{
				UUID:     tt.fields.UUID,
				Name:     tt.fields.Name,
				Cores:    tt.fields.Cores,
				MemoryMb: tt.fields.MemoryMb,
				DiskMb:   tt.fields.DiskMb,
			}
			if err := h.Validate(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Hosting.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewHosting(t *testing.T) {
	type args struct {
		name     string
		cores    int
		memorymb int
		diskmb   int
	}
	tests := []struct {
		name    string
		args    args
		want    *Hosting
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHosting(tt.args.name, tt.args.cores, tt.args.memorymb, tt.args.diskmb)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHosting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHosting() = %v, want %v", got, tt.want)
			}
		})
	}
}
