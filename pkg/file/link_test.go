package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLink(t *testing.T) {
	tmpDir := "link"
	_ = os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	type args struct {
		origin string
		dest   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should link the given dir",
			args: args{
				origin: "testdata",
				dest:   filepath.Join(tmpDir, "testdata"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Link(tt.args.origin, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Link() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !Exists(tt.args.dest) {
				t.Errorf("Link does not work")
			}
		})
	}
}
