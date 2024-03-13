package task

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"go.uber.org/zap"
	"time"
)

func HandleProxyIP() (err error) {

	ipAddr := http.ProxyAddress2DB()

	if ipAddr == "" {
		global.GVA_LOG.Error("获取代理地址失败")
		return fmt.Errorf("")
	} else {
		global.GVA_REDIS.Set(context.Background(), global.SysProxyIPPrefix, ipAddr, 150*time.Second)
		global.GVA_LOG.Info("使用新代理ip，并设置复用池", zap.Any("addr", ipAddr))

	}
	return nil
}
