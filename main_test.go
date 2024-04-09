package main

import "testing"

func Test_msg(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "sucesso", args: struct{ msg string }{msg: "casa"}, want: "casa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := msg(tt.args.msg); got != tt.want {
				t.Errorf("msg() = %v, want %v", got, tt.want)
			}
		})
	}
}
