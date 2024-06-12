package helpers

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/model"
)

func TestUUID_ToUUID(t *testing.T) {
	uuidStr := "bebd6847-f282-4599-9dd9-05a8c0fbfe91"
	uuidVal, _ := uuid.Parse(uuidStr)
	type fields struct {
		Value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    uuid.UUID
		wantErr bool
	}{
		{
			name: "ProtoToUUID transforms model.UUID to uuid.UUID",
			fields: fields{
				Value: uuidStr,
			},
			want:    uuidVal,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &model.UUID{
				Value: tt.fields.Value,
			}
			got, err := ProtoToUUID(id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UUID.ProtoToUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UUID.ProtoToUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID_FromUUID(t *testing.T) {
	uuidStr := "f8fa834a-2bbf-47dc-ac87-e6673a7ed0b3"
	uuidVal, _ := uuid.Parse(uuidStr)
	type args struct {
		uid uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ProtoFromUUID assigns uuid.UUID string value to UUID.Value",
			args: args{
				uid: uuidVal,
			},
			want: uuidStr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := ProtoFromUUID(tt.args.uid)
			if !reflect.DeepEqual(id.Value, tt.want) {
				t.Errorf("UUID.Value = %v, want %v", id.Value, tt.want)
			}
		})
	}
}
