package utils

import (
	"testing"
)

func TestHttp(t *testing.T) {
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
			name:    "GET IP1 test",
			args:    args{"1.2.3.4"},
			wantErr: false,
		},
		{
			name:    "GET IP2 test",
			args:    args{"121.225.97.101"},
			wantErr: false,
		},
		{
			name:    "GET IP3 test",
			args:    args{"0.1.1.0"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region, err := SearchIp2Region(tt.args.d)
			t.Logf("GET() got = %v", region)
			if err != nil {
				t.Errorf("GET() got = %v", err)
			}
		})
	}
}
