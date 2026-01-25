package vtype

import (
	"reflect"
	"testing"
	"time"
)

func TestNewDomainTimeFromString(t *testing.T) {
	type args struct {
		dateStr string
		format  string
	}
	tp := time.Date(2026, 01, 21, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		args    args
		want    *Time
		wantErr bool
	}{
		{
			name: "test_domain_time_for_string",
			args: args{
				dateStr: "2026-01-21",
				format:  "2006-01-02",
			},
			want:    NewTime(&tp),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimeFromString(tt.args.dateStr, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomainTimeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDomainTimeFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDomainTimeFromUnix(t *testing.T) {
	type args struct {
		t int64
	}
	tp := time.Date(2026, 01, 21, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		args args
		want *Time
	}{
		{
			name: "test_domain_time_from_unix",
			args: args{
				t: 1768953600,
			},
			want: NewTime(&tp),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTimeFromUnix(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDomainTimeFromUnix() = %v, want %v", got, tt.want)
			}
		})
	}
}
