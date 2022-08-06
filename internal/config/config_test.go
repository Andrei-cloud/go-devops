package config

import "testing"

func TestReadConfigFile(t *testing.T) {
	type args struct {
		path string
		c    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"srever cfg",
			args{
				"../../cmd/server/server.cfg",
				&ServerConfig{},
			},
		},
		{
			"agent cfg",
			args{
				"../../cmd/agent/agent.cfg",
				&ServerConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReadConfigFile(tt.args.path, tt.args.c)
		})
	}
}
