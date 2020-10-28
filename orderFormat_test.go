package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_ordrFormat(t *testing.T) {
	type args struct {
		f *os.File
	}

	_f1, _ := os.OpenFile("test/nginx_order_1.conf", os.O_RDONLY, 0755)
	_f2, _ := os.OpenFile("test/nginx_order_2.conf", os.O_RDONLY, 0755)
	_f3, _ := os.OpenFile("test/nginx_order_3.conf", os.O_RDONLY, 0755)

	tests := []struct {
		name    string
		args    args
		wantNgs []ngxString
		wantErr bool
	}{
		{
			name: "sample line",
			args: struct{ f *os.File }{f: _f1},
			wantNgs: []ngxString{
				{
					data: "#user  nobody;",
					tab: 0,
				},
				{
					data: "worker_processes 1;",
					tab: 0,
				},
				{
					data: "",
					tab: 0,
				},
				{
					data: "#error_log  logs/error.log;",
					tab: 0,
				},
				{
					data: "#error_log  logs/error.log  notice;",
					tab: 0,
				},
				{
					data: "#error_log  logs/error.log  info;",
					tab: 0,
				},

			},
		},
		{
			name: "contains { in line",
			args: struct{ f *os.File }{f: _f2},
			wantNgs: []ngxString{
				{
					data: "#pid        logs/nginx.pid;",
					tab: 0,
				},
				{
					data: "events {",
					tab: 0,
				},
				{
					data: "worker_connections 1024;",
					tab: 4,
				},
				{
					data: "}",
					tab: 0,
				},
			},
		},
		{
			name: "contains new line in line",
			args: struct{ f *os.File }{f: _f3},
			wantNgs: []ngxString{
				{
					data: "#pid        logs/nginx.pid;",
					tab: 0,
				},
				{
					data: "events {",
					tab: 0,
				},
				{
					data: "worker_connections 1024;",
					tab: 4,
				},
				{
					data: "log_format main '$remote_addr - $remote_user [$time_local] \"$request\" '",
					tab: 4,
				},
				{
					data: "'$status $body_bytes_sent \"$http_referer\" '",
					tab: 4,
				},
				{
					data: "'\"$http_user_agent\" \"$http_x_forwarded_for\"';",
					tab: 4,
				},
				{
					data: "}",
					tab: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNgs, err := ordrFormat(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ordrFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNgs, tt.wantNgs) {
				t.Errorf("ordrFormat() gotNgs = %v, want %v", gotNgs, tt.wantNgs)
			}
		})
	}
}
