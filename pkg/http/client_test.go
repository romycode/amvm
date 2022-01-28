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
		//http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//	if req.URL.String() == "/dist/index.json" {
		//		rw.Header().Add("Content-Type", "application/json")
		//		_, _ = rw.Write([]byte(`[{"version":"v17.4.0","date":"2021-09-15","npm":"7.24.0","v8":"9.3.345.19","uv":"1.42.0","zlib":"1.2.11","openssl":"1.1.1l+quic","modules":"93","lts":"Fermium","security":false},{"version":"v16.13.2","date":"2021-09-20","npm":"7.24.0","v8":"9.3.345.19","uv":"1.42.0","zlib":"1.2.11","openssl":"1.1.1l+quic","modules":"93","lts":false,"security":false},{"version":"v14.18.3","date":"2021-09-25","npm":"7.24.0","v8":"9.3.345.19","uv":"1.42.0","zlib":"1.2.11","openssl":"1.1.1l+quic","modules":"93","lts":"Fermium","security":true}]`))
		//	}
		//
		//	if req.URL.String() == "/dist/v17.4.0/node-v17.4.0-linux-x64.tar.gz" ||
		//		req.URL.String() == "/dist/v16.13.2/node-v16.13.2-linux-x64.tar.gz" ||
		//		req.URL.String() == "/dist/v14.18.3/node-v14.18.3-linux-x64.tar.gz" {
		//		http.ServeFile(rw, req, "testdata/test.tar.gz")
		//	}
		//}),
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
				url:    ts.URL + "/test",
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
			got, err := f.Request(tt.args.method, tt.args.url, tt.args.data)
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
