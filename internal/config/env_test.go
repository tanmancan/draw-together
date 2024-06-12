package config

import (
	"os"
	"testing"
)

func Test_getEnvInt(t *testing.T) {
	testKey := "testEnvInt"
	testVal := "48"
	wantVal := 48

	testKeyErr := "testEnvIntErr"
	testValErr := "$#"

	testFallback := 99
	type args struct {
		key      string
		fallback int
	}
	tests := []struct {
		name     string
		args     args
		want     int
		setUp    func()
		tearDown func()
	}{
		{
			name: "getEnvInt returns an int value if it exists as an ENV variable",
			args: args{
				key:      testKey,
				fallback: testFallback,
			},
			want: wantVal,
			setUp: func() {
				os.Setenv(testKey, testVal)
			},
			tearDown: func() {
				os.Unsetenv(testKey)
			},
		},
		{
			name: "getEnvInt returns the fallback value if it does not exists as an ENV variable",
			args: args{
				key:      testKey,
				fallback: testFallback,
			},
			want: testFallback,
			setUp: func() {
			},
			tearDown: func() {
				os.Unsetenv(testKey)
			},
		},
		{
			name: "getEnvInt returns the fallback value if there was an error parsing the ENV value",
			args: args{
				key:      testKeyErr,
				fallback: testFallback,
			},
			want: testFallback,
			setUp: func() {
				os.Setenv(testKeyErr, testValErr)
			},
			tearDown: func() {
				os.Unsetenv(testKeyErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp()
			}
			defer func() {
				if tt.tearDown != nil {
					tt.tearDown()
				}
			}()
			if got := getEnvInt(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("getEnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvBool(t *testing.T) {
	testKey := "testEnvBool"
	testVal := "false"
	wantVal := false

	testKeyErr := "testEnvBoolErr"
	testValErr := "(*#@)"

	testFallback := true
	type args struct {
		key      string
		fallback bool
	}
	tests := []struct {
		name     string
		args     args
		want     bool
		setUp    func()
		tearDown func()
	}{
		{
			name: "getEnvBool returns an bool value if it exists as an ENV variable",
			args: args{
				key:      testKey,
				fallback: testFallback,
			},
			want: wantVal,
			setUp: func() {
				os.Setenv(testKey, testVal)
			},
			tearDown: func() {
				os.Unsetenv(testKey)
			},
		},
		{
			name: "getEnvBool returns the fallback value if it does not exists as an ENV variable",
			args: args{
				key:      testKey,
				fallback: testFallback,
			},
			want: testFallback,
			setUp: func() {
			},
			tearDown: func() {
				os.Unsetenv(testKey)
			},
		},
		{
			name: "getEnvBool returns the fallback value if there was an error parsing the ENV value",
			args: args{
				key:      testKeyErr,
				fallback: testFallback,
			},
			want: testFallback,
			setUp: func() {
				os.Setenv(testKeyErr, testValErr)
			},
			tearDown: func() {
				os.Unsetenv(testKeyErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp()
			}
			defer func() {
				if tt.tearDown != nil {
					tt.tearDown()
				}
			}()
			if got := getEnvBool(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("getEnvBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvStr(t *testing.T) {
	testKey := "testEnvStr"
	testVal := "test-string-value"
	wantVal := "test-string-value"
	testFallback := "test-fallback-string-value"
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name     string
		args     args
		want     string
		setUp    func()
		tearDown func()
	}{
		{
			name: "getEnvStr returns an string value if it exists as an ENV variable",
			args: args{
				key:      testKey,
				fallback: testFallback,
			},
			want: wantVal,
			setUp: func() {
				os.Setenv(testKey, testVal)
			},
			tearDown: func() {
				os.Unsetenv(testKey)
			},
		},
		{
			name: "getEnvStr returns the fallback value if it does not exists as an ENV variable",
			args: args{
				key:      testKey,
				fallback: testFallback,
			},
			want: testFallback,
			setUp: func() {
			},
			tearDown: func() {
				os.Unsetenv(testKey)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp()
			}
			defer func() {
				if tt.tearDown != nil {
					tt.tearDown()
				}
			}()
			if got := getEnvStr(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("getEnvStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
