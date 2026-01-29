package vtype

import (
	"reflect"
	"testing"
)

func TestNewEmailAddr(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    *EmailAddr
		wantErr bool
	}{
		{
			name: "test_valid_email",
			args: args{
				address: "john.doe@xyz.com",
			},
			want: &EmailAddr{
				address: "john.doe@xyz.com",
			},
			wantErr: false,
		},
		{
			name: "test_invalid_email",
			args: args{
				address: "john.doe@xyz",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmailAddr(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmailAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmailAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailAddr_MarshalJSON(t *testing.T) {
	type fields struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test_email_addr_marshal_json",
			fields: fields{
				address: "john.doe@xyz.com",
			},
			want:    []byte(`"john.doe@xyz.com"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EmailAddr{
				address: tt.fields.address,
			}
			got, err := e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailAddr_UnmarshalJSON(t *testing.T) {
	type fields struct {
		address string
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *EmailAddr
		wantErr bool
	}{
		{
			name: "test_valid_email",
			args: args{
				bytes: []byte(`"john.doe@xyz.com"`),
			},
			want: &EmailAddr{
				address: "john.doe@xyz.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EmailAddr{
				address: tt.fields.address,
			}
			if err := e.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", e, tt.want)
			}
		})
	}
}
