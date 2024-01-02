package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"go.uber.org/zap"
)

func HandleChannelStatisTask() (err error) {

	global.GVA_LOG.Info("用户通道粒度 --> 统计开始")
	err = service.ServiceGroupApp.VboxServiceGroup.CronVboxBdaChIndexD()
	if err != nil {
		global.GVA_LOG.Error("用户通道粒度 --> 统计失败 ", zap.Error(err))
		return
	} else {
		global.GVA_LOG.Info("用户通道粒度 --> 统计成功")
	}
	global.GVA_LOG.Info("用户通道账户粒度 --> 统计开始")
	err = service.ServiceGroupApp.VboxServiceGroup.CronVboxBdaChaccIndexD()
	if err != nil {
		global.GVA_LOG.Error("用户通道账户粒度 --> 统计失败 ", zap.Error(err))
		return
	} else {
		global.GVA_LOG.Info("用户通道账户粒度 --> 统计成功")
	}

	global.GVA_LOG.Info("用户通道店铺粒度 --> 统计开始")
	err = service.ServiceGroupApp.VboxServiceGroup.CronVboxBdaChShopIndexD()
	if err != nil {
		global.GVA_LOG.Error("用户通道店铺粒度 --> 统计失败 ", zap.Error(err))
		return
	} else {
		global.GVA_LOG.Info("用户通道店铺粒度 --> 统计成功")
	}
	return err
}
