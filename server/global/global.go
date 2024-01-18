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
	ChanAccDailyUsed  = "chan_acc_daily_used:ac_id:%s"
	ChanAccDailyLimit = "chan_acc_daily_limit:ac_id:%s"
	ChanAccTotalUsed  = "chan_acc_total_used:ac_id:%s"
	ChanAccTotalLimit = "chan_acc_total_limit:ac_id:%s"
	ChanAccCountUsed  = "chan_acc_count_used:ac_id:%s"
	ChanAccCountLimit = "chan_acc_count_limit:ac_id:%s"

	BalanceNotEnough       = "当前账户余额不足，请及时充值积分后再开启账号，关闭账号ID: %s, 关闭账号： %s"
	AccDailyLimitNotEnough = "当前账户日消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s"
	AccTotalLimitNotEnough = "当前账户总消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s"
	AccCountLimitNotEnough = "当前账户笔数消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s"
	AccQryRecordsEx        = "当前账户查官方记录异常，请核查CK，无法开启账号，ID: %s, 关闭账号： %s"
	AccDelSuccess          = "删除通道账号成功，ID：%v, 通道账号：%s"

	ResourceNotEnough        = "当前库存不足，请及时核查匹配资源剩余情况"
	ResourceShopNotEnough    = "当前库存不足，请及时核查匹配【引导商铺情况】，请求通道：%s，请求金额：%v"
	ResourcePayCodeNotEnough = "当前库存不足，请及时核查匹配【预产码剩余情况】，请求通道：%s，请求金额：%v"
	ResourceAccNotEnough     = "当前库存不足，请及时核查匹配【通道账号情况】，请求通道：%s，请求金额：%v"

	NotifyEx          = "付方回调异常, err: %v， 付方响应信息： %v"
	NotifyHandSuccess = "【候补单】付方回调成功， 付方响应信息： 状态码：%v，数据包：%v"
)

const (
	BloomFilterErrorRate = 0.001
	BloomFilterCapacity  = 100000

	ChanOrgAccFilter    = "vb_accFilter:org_%s:chan_%s"              // 同组织通道下可用账号（过滤器）
	ChanOrgQBAccZSet    = "vb_acc_qb_set:org_%d:chan_%s:money_%v"    // 同组织通道下可用账号
	ChanOrgJ3AccZSet    = "vb_acc_j3_set:org_%d:chan_%s"             // 同组织通道下可用账号（剑三）
	ChanOrgShopAddrZSet = "vb_shop_addr_set:org_%d:chan_%s:money_%d" // 同组织通道下可用店铺地址

	ChanOrgPayCodeLocZSet     = "vb_pay_code_set:org_%d:chan_%s:money_%d:operator_%s:loc_%s" // 同组织通道下可用付款码（取用池）
	ChanOrgPayCodePrefix      = "vb_pay_code_set:org_%d:chan_%s:money_%d:*"                  // 同组织通道下可用付款码（取用池）
	ChanOrgPayCodeMoneyPrefix = "vb_pay_code_set:org_%d:chan_%s:money_*"                     // 同组织通道下可用付款码（取用池）

	PayOrderKey    = "vb_order:%s"
	PayOrderJUCKey = "vb_juc_order:%s"

	OrgChanSet      = "vb_cid_list:org_%d"           // 组织下拥有的产品id
	OrgShopMoneySet = "vb_shop_money:org_%d:chan_%s" // 同组织通道下可用商铺金额

	PcAccWaiting = "vb_acc_waiting_pc:acid_%s" // 预产类-等待开启的账户(冷却中)

	YdQBAccWaiting = "vb_acc_qb_waiting_yd:acid_%s:money_%v" // 引导类-等待开启的账户(冷却中)

	YdJ3AccWaiting = "vb_acc_j3_waiting_yd:acid_%s" // 引导类-等待开启的账户(冷却中)

)

const (
	ProductRecordQBPrefix = "product_record:qb:proxy" // QB查询
)

const (
	SysUserOrgPrefix = "sys_user:%d"
)

// PAcc 相关

const (
	PayAccPrefix = "p_acc_id:" // 商户信息

	PAccPay    = "p_acc_pay"
	PAccCreate = "p_acc_create"
	PAccQuery  = "p_acc_query"
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
	WalletEventOrderCost      = "积分消费[%d], 来自订单[%s]"  // 消费
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
