package file

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	type args struct {
		name string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should write file correctly",
			args: args{
				name: "testdata/prueba.txt",
				data: nil,
			},
			wantErr: false,
		},
		{
			name: "should return error if can't write file correctly",
			args: args{
				name: "tmp/prueba.txt",
				data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Write(tt.args.name, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		if !tt.wantErr {
			err := os.Remove(tt.args.name)
			if err != nil {
				t.Errorf("error cleaning test")
			}
		}
	}
}
