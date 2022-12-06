package dev

import (
	"testing"
)

func TestPortIsInRange(t *testing.T) {
	type args struct {
		pNum   int
		pRange string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{0, ""}, false},
		{"", args{2, "4,6,12-15,18-"}, false},
		{"", args{4, "4,6,12-15,18-"}, true},
		{"", args{18, "4,6,12-15,18-"}, true},
		{"", args{17, "4,6,12-15,18-"}, false},
		{"", args{25, "4,6,12-15,18-"}, true},
		{"", args{4, "2-15"}, true},
		{"", args{2, "4,6,-18,12-15"}, true},
		{"", args{4, "4,6,-18,12-15"}, true},
		{"", args{18, "4,6,-18,12-15"}, true},
		{"", args{17, "4,6,-18,12-15"}, true},
		{"", args{4, "1,6,8-10"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PortIsInRange(tt.args.pNum, tt.args.pRange); got != tt.want {
				t.Errorf("PortIsInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTarget(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name      string
		args      args
		wantIp    string
		wantPorts string
	}{
		{"empty", args{""}, "", ""},
		{"ok", args{"192.168.0.10"}, "192.168.0.10", ""},
		{"ok", args{"192.168.0.10:12"}, "192.168.0.10", "12"},
		{"ok", args{"192.168.0.10:12,13"}, "192.168.0.10", "12,13"},
		{"ok", args{"192.168.0.10:12,13,14-"}, "192.168.0.10", "12,13,14-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIp, gotPorts := parseTarget(tt.args.target)
			if gotIp != tt.wantIp {
				t.Errorf("parseTarget() gotIp = %v, want %v", gotIp, tt.wantIp)
			}
			if gotPorts != tt.wantPorts {
				t.Errorf("parseTarget() gotPorts = %v, want %v", gotPorts, tt.wantPorts)
			}
		})
	}
}
