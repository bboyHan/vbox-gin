package product

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
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
			name:    "GET 11 test",
			args:    args{"975796969"},
			wantErr: false,
		},
		{
			name:    "GET 22 test",
			args:    args{"33713674"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records := Records(tt.args.d)

			global.GVA_LOG = core.Zap()
			global.GVA_LOG.Info("ret:  ->", zap.Any("water list", records.Payments))

			fmt.Printf("ret: %v \n", records.Payments)
		})
	}
}
