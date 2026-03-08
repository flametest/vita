package vtype

import (
	"reflect"
	"testing"
)

func TestIsValidPhoneNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidPhoneNumber",
			args: args{
				number: "+8613512345678",
			},
			want: true,
		},
		{
			name: "InvalidPhoneNumber",
			args: args{
				number: "13512345678",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPhoneNumber(tt.args.number); got != tt.want {
				t.Errorf("IsValidPhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPhoneNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name    string
		args    args
		want    *PhoneNumber
		wantErr bool
	}{
		{
			name: "NewPhoneNumber",
			args: args{
				number: "+8613512345678",
			},
			want: &PhoneNumber{
				countryCode:      "86",
				subscriberNumber: "13512345678",
			},
			wantErr: false,
		},
		{
			name: "NewPhoneNumberWithError",
			args: args{
				number: "13512345678",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPhoneNumber(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPhoneNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhoneNumber_MarshalJSON(t *testing.T) {
	type fields struct {
		countryCode      string
		subscriberNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "MarshalJSON",
			fields: fields{
				countryCode:      "86",
				subscriberNumber: "13512345678",
			},
			want:    []byte(`"+8613512345678"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PhoneNumber{
				countryCode:      tt.fields.countryCode,
				subscriberNumber: tt.fields.subscriberNumber,
			}
			got, err := p.MarshalJSON()
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

func TestPhoneNumber_UnmarshalJSON(t *testing.T) {
	type fields struct {
		countryCode      string
		subscriberNumber string
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UnmarshalJSON",
			fields: fields{
				countryCode:      "86",
				subscriberNumber: "13512345678",
			},
			args: args{
				bytes: []byte(`"+8613512345678"`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PhoneNumber{}
			if err := p.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if p.countryCode != tt.fields.countryCode {
				t.Errorf("UnmarshalJSON() countryCode = %v, want %v", p.countryCode, tt.fields.countryCode)
			}
			if p.subscriberNumber != tt.fields.subscriberNumber {
				t.Errorf("UnmarshalJSON() subscriberNumber = %v, want %v", p.subscriberNumber, tt.fields.subscriberNumber)
			}
		})
	}
}
