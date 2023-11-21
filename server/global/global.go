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

// acc 相关

const (
	ChanAccDailyUsed       = "chan_acc_daily_used:ac_id:%s"
	ChanAccTotalUsed       = "chan_acc_total_used:ac_id:%s"
	ChanAccCountUsed       = "chan_acc_count_used:ac_id:%s"
	BalanceNotEnough       = "当前账户余额不足，请及时充值积分后再开启账号，关闭账号ID: %s, 关闭账号： %s"
	AccDailyLimitNotEnough = "当前账户日消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s"
	AccTotalLimitNotEnough = "当前账户总消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s"
	AccCountLimitNotEnough = "当前账户笔数消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s"
	AccQryRecordsEx        = "当前账户查官方记录异常，请核查CK，无法开启账号，ID: %s, 关闭账号： %s"
)

const (
	BloomFilterErrorRate = 0.001
	BloomFilterCapacity  = 100000

	ChanOrgAccFilter = "vb_accFilter:org_%s:chan_%s"             // 同组织通道下可用账号（过滤器）
	ChanOrgAccZSet   = "vb_acc_set:org_%s:chan_%s"               // 同组织通道下可用账号（过滤器）
	ChanOrgShopZSet  = "vb_acc_set:org_%d:shop_chan_%s:money_%d" // 同组织通道下可用店铺（过滤器）
)

const (
	ProductRecordQBPrefix = "product_record:qb:proxy" // QB查询
)

const (
	UserOrgChannelCodePrefix = "user_org_channel_code_ids:" // 组织下拥有的产品id
	PayAccPrefix             = "pacc_id:"                   // 商户信息
)

// PAcc 相关

const (
	PAccPay    = "pacc_pay"
	PAccCreate = "pacc_create"
	PAccQuery  = "pacc_query"
)

// 用户充值消费相关

const (
	WalletRechargeType = 1 // 直充
	WalletTransferType = 2 // 划转
	WalletOrderType    = 3 // 订单积分消费
)

const (
	WalletEventRechargePrefix = "VBIN"                // 充值
	WalletEventTransferPrefix = "VBTS"                // 划转
	WalletEventOrderPrefix    = "VBOD"                // 订单消费
	WalletEventRecharge       = "充值积分[%d]"            // 充值
	WalletEventTransfer       = "积分扣减[%d], 积分划转至[%s]" // 划转
	WalletEventIncome         = "积分增加[%d], 来自[%s]"    // 划转
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
