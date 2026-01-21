package vo

import (
	"reflect"
	"testing"
)

func TestPassword_MarshalText(t *testing.T) {
	type fields struct {
		hashedPwd string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test Masking",
			fields: fields{
				hashedPwd: "hashedPwd",
			},
			want:    []byte("******"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Password{
				hashedPwd: tt.fields.hashedPwd,
			}
			// got, err := json.Marshal(p)
			got, err := p.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalText() got = %v, want %v", got, tt.want)
			}
		})
	}
}
