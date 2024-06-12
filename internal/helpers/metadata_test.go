package helpers

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func TestMetadataGetCi(t *testing.T) {
	md := metadata.MD{}
	testKey := "TestKey"
	testVal := uuid.NewString()
	testKey2 := "TestKey2"
	testVal2 := uuid.NewString()
	testVal3 := uuid.NewString()
	md.Set(testKey, testVal)
	md.Set(testKey2, testVal2, testVal3)
	type args struct {
		md metadata.MD
		k  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "MetadataGetCi should return value when using original key case",
			args: args{
				md: md,
				k:  testKey,
			},
			want: []string{testVal},
		},
		{
			name: "MetadataGetCi should return value when using ALL_CAPS key",
			args: args{
				md: md,
				k:  "TESTKEY",
			},
			want: []string{testVal},
		},
		{
			name: "MetadataGetCi should return value when using lower case key",
			args: args{
				md: md,
				k:  "testkey",
			},
			want: []string{testVal},
		},
		{
			name: "MetadataGetCi should return nil when using and unknown key",
			args: args{
				md: md,
				k:  "unknown+key",
			},
			want: nil,
		},
		{
			name: "MetadataGetCi should return slice when using key with list value",
			args: args{
				md: md,
				k:  "testkey2",
			},
			want: []string{testVal2, testVal3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MetadataGetCi(tt.args.md, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetadataGetCi() = %v, want %v", got, tt.want)
			}
		})
	}
}
