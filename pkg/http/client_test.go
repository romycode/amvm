package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.URL.String() == "/test" {
				rw.Header().Add("Content-Type", "application/json")
				_, _ = rw.Write([]byte(`{"data":"OK!"}`))
			}
		}),
	)
}

func TestClient_Request(t *testing.T) {
	type fields struct {
		hc  *http.Client
		url string
	}
	type args struct {
		method string
		url    string
		data   string
	}

	ts := testServer()
	defer ts.Close()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "it should send get request",
			fields: fields{
				hc:  ts.Client(),
				url: ts.URL,
			},
			args: args{
				method: "GET",
				url:    "/test",
				data:   "",
			},
			want:    []byte(`{"data":"OK!"}`),
			wantErr: false,
		},
		{
			name: "it should return error if give invalid URL",
			fields: fields{
				hc:  ts.Client(),
				url: ts.URL,
			},
			args: args{
				method: "GET",
				url:    "~://hello.com/test",
				data:   "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewClient(ts.Client(), tt.fields.url)
			got, err := f.Request(tt.args.method, f.URL()+tt.args.url, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				data, _ := io.ReadAll(got.Body)
				if !reflect.DeepEqual(data, tt.want) {
					t.Errorf("Request() got = %v, want %v", data, tt.want)
				}
			}
		})
	}
}
