package main

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_getDelimLine(t *testing.T) {
	type args struct {
		f *os.File
	}

	_f, err := os.OpenFile("test/nginx.conf", os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		wantStr []string
		wantErr bool
	}{
		{
			name:    "sample file",
			args:    struct{ f *os.File }{f: _f},
			wantStr: []string{"#user  nobody;", "\nworker_processes 1;", "\n\n#error_log  logs/error.log;", "\n#error_log  logs/error.log  notice;", "\n#error_log  logs/error.log  info;", "\n\n#pid        logs/nginx.pid;", "\nevents {\n    worker_connections 1024;", "\n}\n\n\nhttp {\n    lua_package_path \"/path/to/lua-resty-upstream-healthcheck/lib/?.lua;", ";", "\";", "\n    include mime.types;", "\n    default_type application/octet-stream;", "\n\n    log_format main '$remote_addr - $remote_user [$time_local] \"$request\" '\n    '$status $body_bytes_sent \"$http_referer\" '\n    '\"$http_user_agent\" \"$http_x_forwarded_for\"';", "\n\n    #access_log  logs/access.log  main;", "\n    sendfile on;", "\n    #tcp_nopush     on;", "\n\n    #keepalive_timeout  0;", "\n    keepalive_timeout 65;", "\n\n    #gzip  on;", "\n\n    # server {\n    #     listen 80;", "\n    #     server_name localhost;", "\n    #     #charset koi8-r;", "\n    #     #access_log  logs/host.access.log  main;", "\n    #     location / {\n    #         root html;", "\n    #         index index.html index.htm;", "\n    #     }\n    #     #error_page  404              /404.html;", "\n    #     # redirect server error pages to the static page /50x.html\n    #     #\n    #     error_page 500 502 503 504 /50x.html;", "\n    #     location = /50x.html {\n    #         root html;", "\n    #     }\n    #     # proxy the PHP scripts to Apache listening on 127.0.0.1:80\n    #     #\n    #     #location ~ \\.php$ {\n    #     #    proxy_pass   http://127.0.0.1;", "\n    #     #}\n    #     # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000\n    #     #\n    #     #location ~ \\.php$ {\n    #     #    root           html;", "\n    #     #    fastcgi_pass   127.0.0.1:9000;", "\n    #     #    fastcgi_index  index.php;", "\n    #     #    fastcgi_param  SCRIPT_FILENAME  /scripts$fastcgi_script_name;", "\n    #     #    include        fastcgi_params;", "\n    #     #}\n    #     # deny access to .htaccess files, if Apache's document root\n    #     # concurs with nginx's one\n    #     #\n    #     #location ~ /\\.ht {\n    #     #    deny  all;", "\n    #     #}\n    # }\n    # another virtual host using mix of IP-, name-, and port-based configuration\n    #\n    #server {\n    #    listen       8000;", "\n    #    listen       somename:8080;", "\n    #    server_name  somename  alias  another.alias;", "\n    #    location / {\n    #        root   html;", "\n    #        index  index.html index.htm;", "\n    #    }\n    #}\n    # HTTPS server\n    #\n    #server {\n    #    listen       443 ssl;", "\n    #    server_name  localhost;", "\n    #    ssl_certificate      cert.pem;", "\n    #    ssl_certificate_key  cert.key;", "\n    #    ssl_session_cache    shared:SSL:1m;", "\n    #    ssl_session_timeout  5m;", "\n    #    ssl_ciphers  HIGH:!aNULL:!MD5;", "\n    #    ssl_prefer_server_ciphers  on;", "\n    #    location / {\n    #        root   html;", "\n    #        index  index.html index.htm;", "\n    #    }\n    #}\n    include conf.d/*.conf;"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, err := getDelimLine(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNextLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStr, tt.wantStr) {
				t.Errorf("getNextLine() gotStr = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func Test_parseStr(t *testing.T) {
	type args struct {
		str   string
		level int
	}
	var tests = []struct {
		name string
		args args
		want []ngxStr
	}{
		{
			name: "mulitple colons",
			args: struct {
				str   string
				level int
			}{str: "lua_package_path \"/path/to/lua-resty-upstream-healthcheck/lib/?.lua;;\";", level: 1},
			want: []ngxStr{
				{
					data:    "lua_package_path \"/path/to/lua-resty-upstream-healthcheck/lib/?.lua;;\";",
					level:   1,
					convert: false,
				},
			},
		},
		{
			name: "one line and contains spaces",
			args: struct {
				str   string
				level int
			}{str: "#user  nobody;\n  \n", level: 1},
			want: []ngxStr{
				{
					data:    "#user  nobody;",
					level:   1,
					convert: false,
				},
				{
					data:    blankLine,
					level:   1,
					convert: false,
				},
				{
					data:    blankLine,
					level:   1,
					convert: false,
				},
			},
		},
		{
			name: "one line",
			args: struct {
				str   string
				level int
			}{str: "#user  nobody;", level: 1},
			want: []ngxStr{{
				data:    "#user  nobody;",
				level:   1,
				convert: false,
			}},
		},
		{
			name: "two lines",
			args: struct {
				str   string
				level int
			}{str: "#user  nobody;\nworker_processes 1;", level: 1},
			want: []ngxStr{
				{
					data:    "#user  nobody;",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_processes 1;",
					level:   1,
					convert: false,
				},
			},
		},
		{
			name: "three lines and tab",
			args: struct {
				str   string
				level int
			}{str: "events {\n    worker_connections 1024;\n}", level: 1},
			want: []ngxStr{
				{
					data:    "events {",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_connections 1024;",
					level:   2,
					convert: false,
				},
				{
					data:    "}",
					level:   1,
					convert: false,
				},
			},
		},
		{
			name: "four lines and two tabs",
			args: struct {
				str   string
				level int
			}{str: "events {\n    worker_connections 1024;\n    worker_connections {\n worker_connections;\n}\n }", level: 1},
			want: []ngxStr{
				{
					data:    "events {",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_connections 1024;",
					level:   2,
					convert: false,
				},
				{
					data:    "worker_connections {",
					level:   2,
					convert: false,
				},
				{
					data:    "worker_connections;",
					level:   3,
					convert: false,
				},
				{
					data:    "}",
					level:   2,
					convert: false,
				},
				{
					data:    "}",
					level:   1,
					convert: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseStr(tt.args.str, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_output(t *testing.T) {
	type args struct {
		ngxs []ngxStr
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sample output",
			args: struct{ ngxs []ngxStr }{ngxs: []ngxStr{
				{
					data:    "#user  nobody;",
					level:   1,
					convert: false,
				},
			}},
			want: "#user  nobody;\n",
		},
		{
			name: "tow outputs",
			args: struct{ ngxs []ngxStr }{ngxs: []ngxStr{
				{
					data:    "#user  nobody;",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_processes 1;",
					level:   1,
					convert: false,
				},
			}},
			want: "#user  nobody;\nworker_processes 1;\n",
		},
		{
			name: "tow outputs",
			args: struct{ ngxs []ngxStr }{ngxs: []ngxStr{
				{
					data:    "#user  nobody;",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_processes 1;",
					level:   1,
					convert: false,
				},
			}},
			want: "#user  nobody;\nworker_processes 1;\n",
		},
		{
			name: "three outputs and tabs",
			args: struct{ ngxs []ngxStr }{ngxs: []ngxStr{
				{
					data:    "events {",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_connections 1024;",
					level:   2,
					convert: false,
				},
				{
					data:    "}",
					level:   1,
					convert: false,
				},
			}},
			want: "events {\n    worker_connections 1024;\n}\n",
		},
		{
			name: "four outputs and tabs",
			args: struct{ ngxs []ngxStr }{ngxs: []ngxStr{
				{
					data:    "events {",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_connections 1024;",
					level:   2,
					convert: false,
				},
				{
					data:    "worker_connections {",
					level:   2,
					convert: false,
				},
				{
					data:    "worker_connections;",
					level:   3,
					convert: false,
				},
				{
					data:    "}",
					level:   2,
					convert: false,
				},
				{
					data:    "}",
					level:   1,
					convert: false,
				},
			}},
			want: "events {\n    worker_connections 1024;\n    worker_connections {\n        worker_connections;\n    }\n}\n",
		},
		{
			name: "one outputs and tabs",
			args: struct{ ngxs []ngxStr }{ngxs: []ngxStr{
				{
					data:    "#user  nobody;",
					level:   1,
					convert: false,
				},
				{
					data:    blankLine,
					level:   1,
					convert: false,
				},
				{
					data:    blankLine,
					level:   1,
					convert: false,
				},
			}},
			want: "#user  nobody;\n\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := output(tt.args.ngxs); got != tt.want {
				t.Errorf("output() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_theSuffixChar(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "sample suffix char",
			args: struct{ str string }{str: "the last char is space   "},
			want: byte('e'),
		},
		{
			name: "all space char",
			args: struct{ str string }{str: "   "},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := theSuffixChar(tt.args.str); got != tt.want {
				t.Errorf("theSuffixChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCustDelimLine(t *testing.T) {
	type args struct {
		f *os.File
	}

	_f, err := os.OpenFile("test/nginx-delim-01.conf", os.O_RDONLY, 0755)
	defer _f.Close()
	if err != nil {
		log.Fatal(err)
	}

	_f2, err := os.OpenFile("test/nginx-delim-02.conf", os.O_RDONLY, 0755)
	defer _f2.Close()
	if err != nil {
		log.Fatal(err)
	}

	_f3, err := os.OpenFile("test/nginx-delim-03.conf", os.O_RDONLY, 0755)
	defer _f2.Close()
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		wantStr []string
		wantErr bool
	}{
		{
			name:    "sample delim",
			args:    struct{ f *os.File }{f: _f},
			wantStr: []string{"http {\nlua_package_path \"/path/to/lua-resty-upstream-healthcheck/lib/?.lua;;\";\ninclude mime.types;\ndefault_type application/octet-stream;\n}"},
			wantErr: false,
		},
		{
			name:    "multiple lines with delim",
			args:    struct{ f *os.File }{f: _f2},
			wantStr: []string{"http {\nlua_package_path \"/path/to/lua-resty-upstream-healthcheck/lib/?.lua;;\";\ninclude mime.types;\ndefault_type application/octet-stream;\nlog_format main '{ \"@timestamp\": \"$time_iso8601\", '\n'\"hostname\": \"$hostname\", '\n'\"domain\": \"$host\",'\n'\"server_name\": \"$server_name\", '\n'\"http_x_forwarded_for\": \"$http_x_forwarded_for\", '\n'\"xes-app\": \"$upstream_http_server\", '\n'\"remote_addr\": \"$remote_addr\", '\n'\"remote_user\": \"$remote_user\", '\n'\"body_bytes_sent\": $body_bytes_sent, '\n'\"request_time\": $request_time, '\n'\"upstream_response_time\": \"$upstream_response_time\", '\n'\"status\": $status, '\n'\"upstream_status\": \"$upstream_status\", '\n'\"connection_requests\": $connection_requests, '\n'\"request\": \"$request\", '\n'\"request_length\": \"$request_length\", '\n'\"request_method\": \"$request_method\", '\n'\"request_body\": \"$request_body\", '\n'\"http_referrer\": \"$http_referer\", '\n'\"http_cookie\": \"$http_cookie\", '\n'\"http_user_agent\": \"$http_user_agent\",'\n'\"url\": \"$host$uri\",'\n'\"upstr_addr\": \"$upstream_addr\",'\n'\"request_trace_id\": \"$request_trace_id\",'\n'\"xes-request-type\": \"$http_xes_request_type\",'\n'\"scheme\": \"$scheme\",'\n'\"traceId\": \"$http_traceid\",'\n'\"rpcId\": \"$http_rpcid\",'\n'\"sw_trace_id\":\"$sent_http_sw_trace_id\"}';\n}"},
			wantErr: false,
		},
		{
			name:    "one line contains {}",
			args:    struct{ f *os.File }{f: _f3},
			wantStr: []string{"events {\nworker_connections 1024;\n}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, err := getCustDelimLine(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCustDelimLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStr, tt.wantStr) {
				t.Errorf("getCustDelimLine() gotStr = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func Test_parseMultipleLine(t *testing.T) {
	type args struct {
		str   []string
		level int
	}
	tests := []struct {
		name string
		args args
		want []ngxStr
	}{
		{
			name: "sample multiple lines",
			args: struct {
				str   []string
				level int
			}{str: []string{
				"events {",
				"worker_connections 1024;",
				"}",
			}, level: 1},
			want: []ngxStr{
				{
					data:    "events {",
					level:   1,
					convert: false,
				},
				{
					data:    "worker_connections 1024;",
					level:   2,
					convert: false,
				},
				{
					data:    "}",
					level:   1,
					convert: false,
				},
			},
		},
		{
			name: "multiple lines",
			args: struct {
				str   []string
				level int
			}{str: []string{"http {", "lua_package_path;", "include mime.types;", "log_format main", "end;", "}"}, level: 1},
			want: []ngxStr{
				{
					data:    "http {",
					level:   1,
					convert: false,
				},
				{
					data:    "lua_package_path;",
					level:   2,
					convert: false,
				},
				{
					data:    "include mime.types;",
					level:   2,
					convert: false,
				},
				{
					data:    "log_format main",
					level:   2,
					convert: false,
				},
				{
					data:    "end;",
					level:   2,
					convert: false,
				},
				{
					data:    "}",
					level:   1,
					convert: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseMultipleLine(tt.args.str, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMultipleLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSingle(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		{
			name:  "multiple lines",
			args:  struct{ s string }{s: "http {\nlua_package_path;\ninclude mime.types;\nlog_format main\nend;\n}"},
			want:  []string{"http {", "lua_package_path;", "include mime.types;", "log_format main", "end;", "}"},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isSingle(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("isSingle() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("isSingle() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
