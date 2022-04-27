package strReplacer

import (
	"testing"
)

func Test_replace_AddParams(t *testing.T) {
	type args struct {
		name    string
		handler func(...string) string
		args    bool
	}
	tests := []struct {
		name string
		r    *Replace
		args args
	}{
		{"名字", Default(), args{"name", func(...string) string { return "相思" }, false}},
		{"性别", Default(), args{"gender", func(...string) string { return "男" }, false}},
		{"网址", Default(), args{"url", func(...string) string { return "http://www.toolv.cn" }, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.AddParams(tt.args.name, tt.args.handler, tt.args.args)
			t.Logf("%+v", tt.r.Params)
		})
	}
}

func Test_replace_AddRegexParams(t *testing.T) {
	type args struct {
		name    string
		handler func([]string, ...string) string
		args    bool
	}
	tests := []struct {
		name string
		r    *Replace
		args args
	}{
		{"名字", Default(), args{`name(\d+)`, func(params []string, args ...string) string {
			return "相思" + params[0]
		}, false}},
		{"性别", Default(), args{`gender(\d+)`, func(params []string, args ...string) string {
			return "男" + params[0]
		}, false}},
		{"网址", Default(), args{`url(\d+)`, func(params []string, args ...string) string {
			if params[0] == "1" {
				return "https://www.toolv.cn"
			}
			return "http://www.toolv.cn"
		}, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.AddRegexParams(tt.args.name, tt.args.handler, tt.args.args)
			t.Logf("%+v", tt.r.RegexParams)
		})
	}
}

func Test_replace_replace(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		r    *Replace
		args args
	}{
		{"名字", Default(), args{"{#name}"}},
		{"名字", Default(), args{"{#name(a,b,c)}"}},
		{"名字", Default(), args{"{#name(1,2,3)}"}},
		{"名字", Default(), args{"{#name(1,2,3,4)}"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.replace(tt.args.s)
			t.Logf("%s", got)
		})
	}
}
