package main

import "testing"

func Test_toPath(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "#1", args: args{arg: "https://github.com/nu50218/gits"}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#2", args: args{arg: "https://github.com/nu50218/gits.git"}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#3", args: args{arg: "github.com/nu50218/gits"}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#4", args: args{arg: "github.com/nu50218/gits.git"}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#5", args: args{arg: `"https://github.com/nu50218/gits"`}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#6", args: args{arg: `"https://github.com/nu50218/gits.git"`}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#7", args: args{arg: `"github.com/nu50218/gits"`}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#8", args: args{arg: `"github.com/nu50218/gits.git"`}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#9", args: args{arg: `"https://github.com/nu50218/gits/"`}, want: "github.com/nu50218/gits", wantErr: false},
		{name: "#10", args: args{arg: `"https://github.com/nu50218/gits.git/"`}, want: "github.com/nu50218/gits", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toDirPath(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("toPath(%s) error = %v, wantErr %v", tt.args.arg, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toPath(%s) = %v, want %v", tt.args.arg, got, tt.want)
			}
		})
	}
}
