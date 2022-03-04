package file

import "testing"

func TestCheck(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "if directory no exist it returns false",
			args: args{name: "no_exist"},
			want: false,
		},
		{
			name: "if directory exist it returns true",
			args: args{name: "testdata"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.args.name); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
