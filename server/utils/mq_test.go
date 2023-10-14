package utils

import (
	"testing"
)

func TestMQ(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		url     string
		wantErr bool
	}{
		{
			name:    "GET IP3 test",
			args:    args{"0.1.1.0"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//MQTask()
		})
	}
}
