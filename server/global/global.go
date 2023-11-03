package global

import (
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"github.com/qiniu/qmgo"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"sync"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	ProductRecordQBPrefix = "product_record:qb:proxy" //QB查询
)

const (
	UserOrgChannelCodePrefix = "user_org_channel_code_ids:" //组织下拥有的产品id
	PayAccPrefix             = "pacc_id:"                   //商户信息
)

const (
	WalletRechargeType = 1 // 直充
	WalletTransferType = 2 // 划转
	WalletOrderType    = 3 // 订单积分消费
)

const (
	WalletEventRechargePrefix = "VBIN"                // 充值或者划转
	WalletEventTransferPrefix = "VBTS"                // 充值或者划转
	WalletEventOrderPrefix    = "VBOD"                // 充值或者划转
	WalletEventRecharge       = "充值积分[%d]"            // 充值或者划转
	WalletEventTransfer       = "积分扣减[%d], 积分划转至[%s]" // 充值或者划转
	WalletEventIncome         = "积分增加[%d], 来自[%s]"    // 充值或者划转
)

var (
	GVA_DB     *gorm.DB
	GVA_DBList map[string]*gorm.DB
	GVA_REDIS  *redis.Client
	GVA_MONGO  *qmgo.QmgoClient
	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	// GVA_LOG    *oplogging.Logger
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}

	BlackCache local_cache.Cache
	lock       sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
