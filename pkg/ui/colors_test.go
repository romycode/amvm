package ui

import "testing"

func TestColorize(t *testing.T) {
	type args struct {
		msg   string
		color Color
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "it should add ui to string",
			args: args{
				msg:   "test string",
				color: Blue,
			},
			want: string(Blue + "test string" + Reset),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Colorize(tt.args.msg, tt.args.color); got != tt.want {
				t.Errorf("Colorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
