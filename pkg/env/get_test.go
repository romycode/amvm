package env

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	_ = os.Setenv("TEST_ENV_VAR", "test-value")

	type args struct {
		name     string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "it should return host environment value",
			args: args{
				name:     "TEST_ENV_VAR",
				fallback: "",
			},
			want: "test-value",
		},
		{
			name: "it should return default value if not exists",
			args: args{
				name:     "NO_EXISTS",
				fallback: "test-value",
			},
			want: "test-value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.name, tt.args.fallback); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
