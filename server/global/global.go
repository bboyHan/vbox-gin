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

	ChanAccOrgCountUnused = "chan_acc_count_unused:org_%d"

	AccQryEx = "系统关号,检测账号CK异常.无法正确查单. ID: %s, acc: %s"

	AccRecord   = "accID,%v"
	UserRecord  = "userID,%v"
	OrderRecord = "orderID,%v"

	OrderStartMsg               = "创建订单，等待用户进行访问..."
	OrderWaitingMsg             = "用户访问订单页，开始进行库存匹配"
	OrderWaitingFinishedMsg     = "库存匹配完成，匹配账号：%s，ID:%s，等待支付..."
	OrderWaitingShopFinishedMsg = "QN库存匹配完成，匹配商品：%s，链接:%s"
	OrderConfirmMsg             = "订单核验充值已到账【%s】，等待发起回调..."
	OrderConfirmErrMsg          = "系统无法正确查单，通过CK查验充值记录异常，账号：%s，ID:%s，进行关闭账号处理，该订单置为失败单"
	OrderTimeoutMsg             = "订单已过期，关闭订单"
	OrderCallbackMsg            = "开始执行回调任务..."
	OrderCallbackRespMsg        = "回调完成，响应状态码：%v，响应数据: %v"
	OrderCallbackFinishedMsg    = "更新回调状态，订单交易完成"
	OrderManualOperationMsg     = "人工核实，进行手动补单入库，订单交易完成"

	OrderConfirmBindMsg           = "订单核验已绑定卡密，卡密信息：%v，等待发起回调..."
	OrderConfirmBindOtherMoneyMsg = "订单核验已绑定卡密，卡密信息：%v，但与订单金额不符，需人工核实"
	OrderConfirmBindQryRetMsg     = "查询卡密合法性，卡密信息：%v，查询结果：%v"
	OrderConfirmBindRetMsg        = "绑定卡密信息：%v，执行结果：%v"
	OrderConfirmBindLimitErrMsg   = "提交的错误次数过多，请重新下单"
	OrderConfirmBindPoolErrMsg    = "查单池CK账户不足，请及时核查"
	OrderConfirmBindErrMsg        = "第【%d】次上传卡密信息，提交卡密信息：%v，执行结果：%v"

	BalanceNotEnough       = "当前账户余额不足，请及时充值积分后再开启账号，关闭账号ID: %s, 关闭账号： %s"
	AccDailyLimitNotEnough = "当前账户日消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s, 当前日消费：%v, 当前限额： %v"
	AccTotalLimitNotEnough = "当前账户总消费已经超出限额，无法开启账号，ID: %s, 关闭账号： %s, 当前总消费：%v, 当前限额： %v"
	AccInCntLimitNotEnough = "当前账户进单数已经超出限额，无法开启账号，ID: %s, 关闭账号： %s, 当前进单数：%v, 当前限额数： %v"
	AccCountLimitNotEnough = "当前账户拉单数已经超出限额，无法开启账号，ID: %s, 关闭账号： %s, 当前拉单数：%v, 当前限额数： %v"
	AccLimitNotEnough      = "当前账户已经超出限额策略，无法使用账号，ID: %s, 关闭账号： %s, 限额信息： %v"
	AccQryRecordsEx        = "当前账户查官方记录异常，请核查CK，无法使用账号，ID: %s, 关闭账号： %s"
	CardAccQryRecordsEx    = "【查单池】当前账户查官方记录异常，请核查CK，无法开启账号，ID: %s, 关闭账号： %s"
	AccQryJ3RecordsEx      = "当前账户查官方记录异常，请核查报文链接，无法使用账号，ID: %s, 关闭账号： %s"
	AccQryShopEx           = "当前组织需开启至少一个商铺地址，请核查商铺信息【通道ID: %s】，无法开启账号，ID: %s, 关闭账号： %s"
	AccDelSuccess          = "删除通道账号成功，ID：%v, 通道账号：%s"

	ResourceNotEnough        = "当前库存不足，请及时核查匹配资源剩余情况"
	ResourceShopNotEnough    = "当前库存不足，请及时核查匹配【引导商铺情况】，请求通道：%s，请求金额：%v"
	ResourcePayCodeNotEnough = "当前库存不足，请及时核查匹配【预产码剩余情况】，请求通道：%s，请求金额：%v"
	ResourceAccNotEnough     = "当前库存不足，请及时核查匹配【通道账号情况】，请求通道：%s，请求金额：%v"
	ResourceQNShopNotEnough  = "当前库存不足，请及时核查匹配【通道商品情况】，请求通道：%s，请求金额：%v"

	NotifyEx          = "付方回调异常, err: %v， 付方响应信息： %v"
	NotifyHandSuccess = "【候补单】付方回调成功， 付方响应信息： 状态码：%v，数据包：%v"
)

const (
	BloomFilterErrorRate = 0.001
	BloomFilterCapacity  = 100000

	MsgFilterKey            = "vb_msg_filter_set"                     // 同组织通道下可用账号（过滤器）
	MsgFilterMem            = "vb_msg_filter:msgID_%s:orderID_%s"     // 同组织通道下可用账号（过滤器）
	ChanOrgProdMoneyAccZSet = "vb_acc_%s_set:org_%d:chan_%s:money_%v" // 同组织通道下可用账号
	ChanOrgProdAccZSet      = "vb_acc_%s_set:org_%d:chan_%s"          // 同组织通道下可用账号
	YdProdMoneyAccWaiting   = "vb_acc_%s_waiting_yd:acid_%s:money_%v" // 引导类-等待开启的账户(冷却中)
	YdProdAccWaiting        = "vb_acc_%s_waiting_yd:acid_%s"          // 引导类-等待开启的账户(冷却中)

	ChanOrgAccFilter        = "vb_accFilter:org_%s:chan_%s"              // 同组织通道下可用账号（过滤器）
	ChanOrgQBAccZSet        = "vb_acc_qb_set:org_%d:chan_%s:money_%v"    // 同组织通道下可用账号
	ChanOrgDnfAccZSet       = "vb_acc_dnf_set:org_%d:chan_%s:money_%v"   // 同组织通道下可用账号
	ChanOrgSdoAccZSet       = "vb_acc_sdo_set:org_%d:chan_%s:money_%v"   // 同组织通道下可用账号
	ChanOrgQBAccZSetPrefix  = "vb_acc_qb_set:org_%d:chan_%s:money_*"     // 同组织通道下可用账号
	ChanOrgDnfAccZSetPrefix = "vb_acc_dnf_set:org_%d:chan_%s:money_*"    // 同组织通道下可用账号
	ChanOrgSdoAccZSetPrefix = "vb_acc_sdo_set:org_%d:chan_%s:money_*"    // 同组织通道下可用账号
	ChanOrgJ3AccZSet        = "vb_acc_j3_set:org_%d:chan_%s"             // 同组织通道下可用账号（剑三）
	ChanOrgQNAccZSet        = "vb_acc_qn_set:org_%d:chan_%s"             // 同组织通道下可用账号（qn）
	ChanOrgECAccZSet        = "vb_acc_ec_set:org_%d:chan_%s"             // 同组织通道下可用账号（ec）
	ChanOrgECPoolAccZSet    = "vb_acc_ec_checkPool_set:org_%d:chan_%s"   // 同组织通道下查单池（ec）
	ChanOrgShopAddrZSet     = "vb_shop_addr_set:org_%d:chan_%s:money_%d" // 同组织通道下可用店铺地址

	ChanOrgQNShopZSet = "vb_qn_shop_set:org_%d:chan_%s:money_%v" // 同组织通道下可用账号（qn acc）

	ChanOrgPayCodeLocZSet     = "vb_pay_code_set:org_%d:chan_%s:money_%d:operator_%s:loc_%s" // 同组织通道下可用付款码（取用池）
	ChanOrgPayCodePrefix      = "vb_pay_code_set:org_%d:chan_%s:money_%d:*"                  // 同组织通道下可用付款码（取用池）
	ChanOrgPayCodeMoneyPrefix = "vb_pay_code_set:org_%d:chan_%s:money_*"                     // 同组织通道下可用付款码（取用池）

	ProdCodeKey = "vb_prod_code:%v"
	ProdTypeKey = "vb_prod_type:%v"
	ProductKey  = "vb_prod_info:%v"

	PayQNShopMoneyKey   = "vb_qn_shop_id:%s:%v"
	PayAccMoneyKey      = "vb_ac_id:%s:%v"
	PayAccKey           = "vb_ac_id:%s"
	PayOrderKey         = "vb_order:%s"
	PayOrderExtLimitKey = "vb_order_ext_limit:%s"
	PayOrderJUCKey      = "vb_juc_order:%s"

	OrgChanSet      = "vb_cid_list:org_%d"           // 组织下拥有的产品id
	OrgShopMoneySet = "vb_shop_money:org_%d:chan_%s" // 同组织通道下可用商铺金额

	PcAccWaiting     = "vb_acc_waiting_pc:acid_%s"  // 预产类-等待开启的账户(冷却中)
	J3AccBalanceZSet = "vb_acc_j3_balance:ac_id:%s" // 剑三账户余额

	ProdAccMoneyWaiting = "vb_acc_%s_waiting_yd:acid_%s:money_%v" // 引导类-等待开启的账户(冷却中)
	ProdAccWaiting      = "vb_acc_%s_waiting_yd:acid_%s"          // 引导类-等待开启的账户(冷却中)

	YdQBAccWaiting  = "vb_acc_qb_waiting_yd:acid_%s:money_%v"  // 引导类-等待开启的账户(冷却中)
	YdDnfAccWaiting = "vb_acc_dnf_waiting_yd:acid_%s:money_%v" // 引导类-等待开启的账户(冷却中)
	YdSdoAccWaiting = "vb_acc_sdo_waiting_yd:acid_%s:money_%v" // 引导类-等待开启的账户(冷却中)

	YdQNAccWaiting  = "vb_acc_qn_waiting_yd:acid_%s"            // 引导类-等待开启的账户(冷却中)
	YdQNShopWaiting = "vb_shop_qn_waiting_yd:shop_mid_%s:id_%v" // 引导类-等待开启的账户(冷却中)

	YdJ3AccWaiting     = "vb_acc_j3_waiting_yd:acid_%s"      // 引导类-等待开启的账户(冷却中)
	YdECAccWaiting     = "vb_acc_ec_waiting_yd:acid_%s"      // 引导类-等待开启的账户(冷却中)
	YdECPoolAccWaiting = "vb_acc_ec_pool_waiting_yd:acid_%s" // 引导类-等待开启的账户(冷却中)

	YdECJdCodeZSet = "vb_acc_ec_jd_code" // 引导类-等待开启的账户(冷却中)

	PayOrderVOKey = "vb_order_vo:%s"
)

const (
	ProductRecordQBPrefix  = "product_record:qb:proxy"  // QB查询
	ProductRecordJ3Prefix  = "product_record:j3:proxy"  // QB查询
	ProductRecordSdoPrefix = "product_record:sdo:proxy" // sdo 查询
	ProductRecordQNPrefix  = "product_record:qn:proxy"  // sdo 查询
)

const (
	SysUserOrgPrefix = "sys_user:%d"

	SysProxyIPPrefix = "sys_proxy"
	SysProxyIPLimit  = "sys_proxy_limit"
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
	UserType  = 1 // 用户级别
	AccType   = 2 // 账户级别
	OrderType = 3 // 订单级别
)

const (
	WalletEventRechargePrefix = "VBIN"                           // 充值
	WalletEventTransferPrefix = "VBTS"                           // 划转
	WalletEventOrderPrefix    = "VBOD"                           // 订单消费
	WalletEventRecharge       = "充值积分[%d]"                       // 充值
	WalletEventTransfer       = "积分扣减[%d], 积分划转至[%s]"            // 划转
	WalletEventIncome         = "积分增加[%d], 来自用户:[%s]"            // 划转
	WalletEventOrderCost      = "积分消费[%d], 来自(通道:[%s]), 订单:[%s]" // 消费
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
