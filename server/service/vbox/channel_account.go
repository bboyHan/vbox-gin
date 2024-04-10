package vbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	prod "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	http2 "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/flipped-aurora/gin-vue-admin/server/vbUtil"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ChannelAccountService struct {
}

/*
    创建布隆过滤器
	bloomFilterKey := "myFilter"
	bloomFilterErrorRate := 0.01
	bloomFilterCapacity := 10000
	global.GVA_REDIS.Do(context.Background(), "BF.RESERVE", bloomFilterKey, bloomFilterErrorRate, bloomFilterCapacity)
	if err != nil {
		fmt.Println("创建布隆过滤器失败:", err)
		return
	}

	// 将元素添加到布隆过滤器中
	err = global.GVA_REDIS.Do(context.Background(), "BF.ADD", "myFilter", "hello").Err()
	if err != nil {
		panic(err)
	}

	// 检查元素是否存在于布隆过滤器中
	exists, err := global.GVA_REDIS.Do(context.Background(), "BF.EXISTS", "myFilter", "hello").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(exists) // 输出 true

	if boolValue, ok := exists.(bool); ok {
		// 如果接口类型的值是bool类型，那么boolValue就是该值的bool表示
		// 在这里你可以使用boolValue
		fmt.Println("The interface value is a bool:", boolValue)
	} else {
		// 如果接口类型的值不是bool类型，这里的代码将会执行
		fmt.Println("The interface value is not a bool")
	}
	// 尝试检查不存在的元素
	exists, err = global.GVA_REDIS.Do(context.Background(), "BF.EXISTS", "myFilter", "world").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(exists) // 输出 false

	// 清除布隆过滤器
	err = global.GVA_REDIS.Do(context.Background(), "DEL", "myFilter").Err()
	if err != nil {
		panic(err)
	}
*/

func (vcaService *ChannelAccountService) LoginQNByQRCode() (res interface{}, err error) {
	ResMap, err := utils.GetQrImg()
	if err != nil {
		global.GVA_LOG.Error("获取二维码失败:", zap.Error(err))
		return
	}

	//qrImg := ResMap["QrImg"]
	//fmt.Println("登录二维码 -> Base64编码:  " + qrImg)
	//login, err := utils.LoopIfLogin()

	return ResMap, err
}

func (vcaService *ChannelAccountService) LoginQNQrStatusCheck(qrSig string) (ret interface{}, err error) {
	return
}

func (vcaService *ChannelAccountService) LoginQQByQRCode() (res interface{}, err error) {
	ResMap, err := utils.GetQrImg()
	if err != nil {
		global.GVA_LOG.Error("获取二维码失败:", zap.Error(err))
		return
	}

	//qrImg := ResMap["QrImg"]
	//fmt.Println("登录二维码 -> Base64编码:  " + qrImg)
	//login, err := utils.LoopIfLogin()

	return ResMap, err
}

func (vcaService *ChannelAccountService) LoginQQQrStatusCheck(qrSig string) (ret interface{}, err error) {
	LoginSig, err := utils.GetLoginSig()
	if err != nil {
		return nil, err
	}

	QrToken := utils.GetQrToken(qrSig)

	str, err := utils.IfLogin(QrToken, LoginSig, qrSig)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(str, "") {
		return nil, errors.New("未知错误 Line 70，请刷新重试！")
	}

	var res = make(map[string]string)

	s := strings.Split(strings.ReplaceAll(str[strings.Index(str, "(")+1:len(str)-1], "'", ""), ",")
	// 65 二维码已失效 66 二维码未失效 67 已扫描,但还未点击确认 0  已经点击确认,并登录成功
	switch s[0] {
	case "65":
		fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "二维码失效 -> 已重新生成")
		return nil, errors.New("二维码已失效 -> 请重新获取")
	case "66":
		fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "登录二维码获取成功 -> 等待扫描")
		return nil, errors.New("登录二维码获取成功 -> 等待扫描")
	case "67":
		fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "已扫描 -> 请点击允许登录")
		return nil, errors.New("已扫描 -> 请点击允许登录")
	case "0":
		// 已经点击确认,并登录成功
		res["NickName"] = s[5]
		res["Location"] = s[2]
		res["QQ"] = utils.ExtractParamValue(res["Location"], "uin")
		fmt.Println("已经点击确认,并登录成功, QQ号:"+res["QQ"], "QQ昵称:"+res["NickName"])
	default:
		return nil, errors.New("未知错误 Line 104，请刷新重试！")
	}

	Res, err := utils.Credential(res["Location"])
	if err != nil {
		return nil, err
	}

	result, err := utils.GetResult(Res)

	if err != nil {
		return nil, err
	}
	result["qq"] = res["QQ"]

	return result, nil
}

// QueryOrgAccAvailable 查询通道账号的官方记录
func (vcaService *ChannelAccountService) QueryOrgAccAvailable(vca *vbox.ChannelAccount) (res interface{}, err error) {

	// filterKey设置

	// 当前用户所属部门
	allOrgIDs := utils2.GetAllOrgID()

	for _, orgID := range allOrgIDs {
		var accList = make(map[string][]interface{})

		// 当前部门拥有的产品
		chanCodeIDs := utils2.GetChannelCodeByOrgID(orgID)

		// 当前部门的所有用户
		userIDs := utils2.GetUsersByOrgIds([]uint{orgID})

		fmt.Printf("org id : %v  ", orgID)
		fmt.Printf("code ids : %v  ", chanCodeIDs)
		fmt.Printf("user ids : %v\n", userIDs)
		var accTempList []vbox.ChannelAccount

		err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("created_by in ?", userIDs).
			Where("cid in ?", chanCodeIDs).Where("status = ? and sys_status = ?", 1, 1).Find(&accTempList).Error
		if err != nil {
			global.GVA_LOG.Error("查询可用通道账号失败:", zap.Error(err))
			return
		}

		// 根据 chan code 分组
		for _, accT := range accTempList {
			if _, ok := accList[accT.Cid]; ok {
				accList[accT.Cid] = append(accList[accT.Cid], accT.AcAccount)
			} else {
				accList[accT.Cid] = []interface{}{accT.AcAccount}
			}
		}

		// 遍历分组后的acc
		for k, accL := range accList {
			fmt.Printf("key : %s", k)
			fmt.Printf("accL : %v\n", accL)

			blKey := fmt.Sprintf(global.ChanOrgAccFilter, strconv.FormatUint(uint64(orgID), 10), k)
			global.GVA_REDIS.Do(context.Background(), "BF.RESERVE", blKey, global.BloomFilterErrorRate, global.BloomFilterCapacity)

			err = global.GVA_REDIS.Do(context.Background(), "BF.MADD", blKey, "accL", "a", "b").Err()
			if err != nil {
				global.GVA_LOG.Error("查询可用通道账号添加到BL失败:", zap.Error(err))
				return
			}

			var count interface{}
			count, err = global.GVA_REDIS.Do(context.Background(), "BF.COUNT", blKey).Result()
			if err != nil {
				global.GVA_LOG.Error("查询可用通道账号添加到BL失败:", zap.Error(err))
				return
			}

			fmt.Printf("accL count : %v\n", count)

		}

	}

	// deepUserIDs := utils2.GetDeepUserIDs(userId)

	//

	//global.GVA_REDIS.Do(context.Background(), "BF.RESERVE", , global.BloomFilterErrorRate, global.BloomFilterCapacity)
	//if err != nil {
	//	global.GVA_LOG.Error("创建布隆过滤器失败:", zap.Error(err))
	//	return
	//}

	return nil, err
}

// QueryAccOrderHis 查询通道账号的官方记录
func (vcaService *ChannelAccountService) QueryAccOrderHis(vca *vbox.ChannelAccount) (res interface{}, err error) {

	var urlQ, channelCode string

	if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) || global.DnfContains(vca.Cid) { // tx系
		c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordQBPrefix).Result()
		if c == 0 {
			channelCode = "qb_proxy"

			err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and chan= ?", 1, channelCode).
				First(&urlQ).Error

			if err != nil {
				return nil, errors.New("该信道无资源配置")
			}

			global.GVA_REDIS.Set(context.Background(), global.ProductRecordQBPrefix, urlQ, 10*time.Minute)
		} else {
			urlQ, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordQBPrefix).Result()
		}
	} else if global.SdoContains(vca.Cid) {
		c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordSdoPrefix).Result()
		if c == 0 {
			channelCode = "sdo_proxy"

			err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and chan= ?", 1, channelCode).
				First(&urlQ).Error

			if err != nil {
				return nil, errors.New("该信道无资源配置")
			}

			global.GVA_REDIS.Set(context.Background(), global.ProductRecordSdoPrefix, urlQ, 10*time.Minute)
		} else {
			urlQ, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordSdoPrefix).Result()
		}
	} else if global.J3Contains(vca.Cid) {
		c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordJ3Prefix).Result()
		if c == 0 {
			channelCode = "j3_proxy"

			err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and chan=?", 1, channelCode).
				First(&urlQ).Error

			if err != nil {
				return nil, errors.New("该信道无资源配置")
			}

			global.GVA_REDIS.Set(context.Background(), global.ProductRecordJ3Prefix, urlQ, 10*time.Minute)
		} else {
			urlQ, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordJ3Prefix).Result()
		}

	} else if global.DyContains(vca.Cid) {
		c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordDYPrefix).Result()
		if c == 0 {
			channelCode = "db_proxy"

			err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and chan=?", 1, channelCode).
				First(&urlQ).Error

			if err != nil {
				return nil, errors.New("该信道无资源配置")
			}

			global.GVA_REDIS.Set(context.Background(), global.ProductRecordDYPrefix, urlQ, 10*time.Minute)
		} else {
			urlQ, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordDYPrefix).Result()
		}

	} else if global.QNContains(vca.Cid) {
		c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordQNPrefix).Result()
		if c == 0 {
			channelCode = "qn_proxy"

			err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and chan=?", 1, channelCode).
				First(&urlQ).Error

			if err != nil {
				return nil, errors.New("该信道无资源配置")
			}

			global.GVA_REDIS.Set(context.Background(), global.ProductRecordQNPrefix, urlQ, 10*time.Minute)

		} else {
			urlQ, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordQNPrefix).Result()

		}
	}

	if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) || global.DnfContains(vca.Cid) { // tx系

		openID, openKey, err := product.Secret(vca.Token)
		if err != nil {
			return nil, err
		}
		records := product.Records(urlQ, openID, openKey, 24*30*time.Hour)
		if records == nil {
			return nil, fmt.Errorf("该账号ck存在异常，请核查")
		}

		if records.Ret != 0 {
			return nil, fmt.Errorf("该账号ck存在异常，请核查")
		}
		//classifier := product.Classifier(records.WaterList)
		return records, nil
	} else if global.J3Contains(vca.Cid) {

		record, err := product.QryJ3Record(*vca)
		if err != nil {
			return nil, err
		}

		timeMax := strconv.FormatInt(time.Now().Unix(), 10)
		global.GVA_LOG.Info("", zap.Any("timeMax", timeMax))
		accKey := fmt.Sprintf(global.J3AccBalanceZSet, vca.AcId)

		list := global.GVA_REDIS.ZRevRangeByScore(context.Background(), accKey, &redis.ZRangeBy{
			Min:    "0",
			Max:    timeMax,
			Offset: 0,
			Count:  100,
		}).Val()
		var accRecords []prod.J3AccountRecord
		for _, mem := range list {
			// 原格式 keyMem := fmt.Sprintf("%s_%s_%v_%d_%d_%d_%d", v.Obj.OrderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
			keyMem := strings.Split(mem, ",")
			money, _ := strconv.Atoi(keyMem[2])

			accRecord := prod.J3AccountRecord{
				OrderID:    keyMem[0],
				AcAccount:  keyMem[1],
				Money:      money,
				NowTime:    keyMem[3],
				HisBalance: keyMem[4],
				CheckTime:  keyMem[5],
				NowBalance: keyMem[6],
			}

			accRecords = append(accRecords, accRecord)
		}
		ret := &prod.J3Records{
			J3BalanceData: *record,
			List:          accRecords,
		}
		return ret, nil
	} else if global.DyContains(vca.Cid) {

		record, err := product.QryDyRecord(vca.Token)
		if err != nil {
			return nil, err
		}

		timeMax := strconv.FormatInt(time.Now().Unix(), 10)
		global.GVA_LOG.Info("", zap.Any("timeMax", timeMax))
		accKey := fmt.Sprintf(global.DyAccBalanceZSet, vca.AcId)

		list := global.GVA_REDIS.ZRevRangeByScore(context.Background(), accKey, &redis.ZRangeBy{
			Min:    "0",
			Max:    timeMax,
			Offset: 0,
			Count:  100,
		}).Val()
		var accRecords []prod.DyAccountRecord
		for _, mem := range list {
			// 原格式 keyMem := fmt.Sprintf("%s_%s_%v_%d_%d_%d_%d", v.Obj.OrderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
			keyMem := strings.Split(mem, ",")
			money, _ := strconv.Atoi(keyMem[2])

			accRecord := prod.DyAccountRecord{
				OrderID:    keyMem[0],
				AcAccount:  keyMem[1],
				Money:      money,
				NowTime:    keyMem[3],
				HisBalance: keyMem[4],
				CheckTime:  keyMem[5],
				NowBalance: keyMem[6],
			}

			accRecords = append(accRecords, accRecord)
		}
		ret := &prod.DYRecords{
			DyWalletInfoRecord: *record,
			List:               accRecords,
		}
		return ret, nil
	} else if global.SdoContains(vca.Cid) {
		records, err := product.QrySdoDaoYuRecords(*vca)
		if err != nil {
			return nil, err
		}
		return records, nil
	} else if global.QNContains(vca.Cid) {
		endTime := time.Now()
		startTime := endTime.Add(-time.Hour * 24 * 7)
		records, err := product.QryQNRecords(*vca, startTime, endTime, "")
		if err != nil {
			return nil, fmt.Errorf("ck失效，请重新上传, err: " + err.Error())
		}
		return records, nil
	}

	return res, err
}

// CountAcc 查询可用通道的 当前等待取用的账号个数
func (vcaService *ChannelAccountService) CountAcc(orgIds []uint) (ret []vboxResp.OrgUnusedData, err error) {

	for _, orgId := range orgIds {
		ids := utils2.GetUsersByOrgId(orgId)
		var res []vboxResp.ChannelAccountUnused
		err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Select("count(1) as total, cid").
			Where("status = ? and sys_status = ? and created_by in (?)", 1, 1, ids).
			Group("cid").Order("id desc").Find(&res).Error

		for i := range res {
			v := &res[i]
			cid := v.Cid
			var moneyList []string
			userIDs := utils2.GetUsersByOrgId(orgId)

			if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
				Where("cid = ? and created_by in ?", cid, userIDs).Scan(&moneyList).Error; err != nil {
				global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
			}

			prodInfo, _ := vbUtil.GetProductByCode(cid)

			switch {
			case strings.Contains(prodInfo.Ext, "money"):
				var accQueueList []vboxResp.AccQueue
				for _, money := range moneyList {
					//cntKey := fmt.Sprintf(global.ChanOrgProdMoneyAccZSet, prodInfo.ProductId, orgId, cid, money)
					cntKey := fmt.Sprintf(prodInfo.Ext, orgId, cid, money)
					cnt := global.GVA_REDIS.ZCount(context.Background(), cntKey, "0", "0").Val()
					accQueue := vboxResp.AccQueue{
						Money:  money,
						Unused: cnt,
					}
					accQueueList = append(accQueueList, accQueue)
				}
				v.List = accQueueList
			default:
				var accQueueList []vboxResp.AccQueue
				//cntKey := fmt.Sprintf(global.ChanOrgProdAccZSet, prodInfo.ProductId, orgId, cid)
				cntKey := fmt.Sprintf(prodInfo.Ext, orgId, cid)
				cnt := global.GVA_REDIS.ZCount(context.Background(), cntKey, "0", "0").Val()
				accQueue := vboxResp.AccQueue{
					Money:  "default",
					Unused: cnt,
				}
				accQueueList = append(accQueueList, accQueue)
				v.List = accQueueList
			}
		}

		ele := vboxResp.OrgUnusedData{
			OrgId: orgId,
			List:  res,
		}
		ret = append(ret, ele)

	}
	global.GVA_LOG.Info("当前等待取用的账号个数", zap.Any("ret", ret))

	return ret, err
}

// TransferChannelForAcc 账号通道转移
func (vcaService *ChannelAccountService) TransferChannelForAcc(vca *vbox.ChannelAccount, c *gin.Context) (err error) {
	var accDB vbox.ChannelAccount
	if err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).First(&accDB).Error; err != nil {
		global.GVA_LOG.Error("查询失败", zap.Error(err))
		return
	}
	sourceCid := accDB.Cid
	if sourceCid == vca.Cid {
		return fmt.Errorf("已在当前通道【%s】中，无需转移", sourceCid)
	}
	var flag bool
	if global.TxContains(vca.Cid) && global.TxContains(sourceCid) {
		flag = true
	} else if global.DnfContains(vca.Cid) && global.DnfContains(sourceCid) {
		flag = true
	} else if global.SdoContains(vca.Cid) && global.SdoContains(sourceCid) {
		flag = true
	} else if global.J3Contains(vca.Cid) && global.J3Contains(sourceCid) {
		flag = true
	} else if global.PcContains(vca.Cid) && global.PcContains(sourceCid) {
		flag = true
	}
	if !flag {
		return fmt.Errorf("不支持的通道转移操作")
	}

	// 发一条清理旧通道资源的
	go func() {
		conn, err := mq.MQ.ConnPool.GetConnection()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)

		ch, err := conn.Channel()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		body := http2.DoGinContextBody(c)
		oldAccTmp := accDB
		oldAccTmp.Status = 0
		oc := vboxReq.ChanAccAndCtx{
			Obj: oldAccTmp,
			Ctx: vboxReq.Context{
				Body:      string(body),
				ClientIP:  c.ClientIP(),
				Method:    c.Request.Method,
				UrlPath:   c.Request.URL.Path,
				UserAgent: c.Request.UserAgent(),
				UserID:    int(oldAccTmp.CreatedBy),
			},
		}
		marshal, err := json.Marshal(oc)

		err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)

		global.GVA_LOG.Info("转移通道清理旧资源", zap.Any("旧cid", oldAccTmp.Cid), zap.Any("新cid", vca.Cid), zap.Any("转移账号", accDB.AcAccount))
	}()

	err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).Update("cid", vca.Cid).Error
	if err != nil {
		return err
	}

	// 新资源发起资源开启，如果账号状态为开启
	if accDB.Status == 1 {
		go func() {
			var vcaDB vbox.ChannelAccount
			err = global.GVA_DB.Model(vbox.ChannelAccount{}).Where("id =?", vca.ID).First(&vcaDB).Error

			conn, err := mq.MQ.ConnPool.GetConnection()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)

			ch, err := conn.Channel()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
			}

			body := http2.DoGinContextBody(c)

			oc := vboxReq.ChanAccAndCtx{
				Obj: vcaDB,
				Ctx: vboxReq.Context{
					Body:      string(body),
					ClientIP:  c.ClientIP(),
					Method:    c.Request.Method,
					UrlPath:   c.Request.URL.Path,
					UserAgent: c.Request.UserAgent(),
					UserID:    int(vcaDB.CreatedBy),
				},
			}
			marshal, err := json.Marshal(oc)

			global.GVA_LOG.Info("转移通道创建新资源,转移时状态开启了,所以先清了旧的,再开新的", zap.Any("新cid", vcaDB.Cid), zap.Any("转移账号", accDB.AcAccount))

			err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)
		}()
	}

	return err
}

// CreateChannelAccount 创建通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) CreateChannelAccount(vca *vbox.ChannelAccount, c *gin.Context) (err error) {
	vca.AcId = rand_string.RandomInt(8)
	token := vca.Token
	//增加校验
	if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) || global.DnfContains(vca.Cid) {
		_, _, errX := product.Secret(token)
		if errX != nil {
			return errX
		}
		isNum := utils.IsNumeric(vca.AcAccount)
		if !isNum {
			return errors.New("QQ账号输入不合法")
		}
	} else if global.QNContains(vca.Cid) {
		isCK := http2.IsValidCookie(token)
		if !isCK {
			return errors.New("ck信息不合法")
		}
		account, errX := product.FindAccNick(token)
		if errX != nil {
			global.GVA_LOG.Warn("未能解析到nick lid的值， 进行随机赋值")
		}
		vca.AcAccount = account

	} else if global.SdoContains(vca.Cid) {
		isCK := http2.IsValidCookie(token)
		if !isCK {
			return errors.New("ck信息不合法")
		}
	} else if global.J3Contains(vca.Cid) {
		parsedURL, errX := url.Parse(token)
		if errX != nil {
			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
			return errors.New("无效的URL")
		}

		query := parsedURL.Query()
		account := query.Get("account")
		zoneCode := query.Get("zoneCode")
		SN := query.Get("SN")
		sign := query.Get("sign")
		if account == "" || zoneCode == "" || SN == "" || sign == "" {
			return errors.New("账号信息不完整")
		}
		if vca.Cid == "2001" && zoneCode != "z22" {
			return errors.New("仅支持双线二区参数，请核查")
		}
		if vca.Cid == "2002" && zoneCode != "z05" {
			return errors.New("仅支持电信五区参数，请核查")
		}
		vca.AcAccount = account
	} else if global.DyContains(vca.Cid) {
		isCK := http2.IsValidCookie(token)
		if !isCK {
			return errors.New("ck信息不合法")
		}
	} else if global.JymContains(vca.Cid) {
		isUrl := http2.IsJymUrl(token)
		if !isUrl {
			return errors.New("url不合法")
		}
	} else if global.ECContains(vca.Cid) {
		isCK := http2.IsValidCookie(token)
		if !isCK {
			return errors.New("ck信息不合法")
		}
		if vca.AcAccount == "" {
			pin := http2.ParseCookie(token, "pin")
			if pin != "" {
				vca.AcAccount = pin
			}
		}
	} else {
		return errors.New("该信道暂不支持创建账号")
	}

	if vca.AcAccount == "" {
		return errors.New("账号信息不完整")
	} else {
		var count int64
		if err = global.GVA_DB.Debug().Model(&vbox.ChannelAccount{}).Where("ac_account = ? and cid = ?", vca.AcAccount, vca.Cid).Count(&count).Error; err != nil {
			global.GVA_LOG.Warn("系统查询异常", zap.Error(err))
			return errors.New("系统查询异常，请重试或联系管理员")
		}
		if count > 0 {
			global.GVA_LOG.Warn("账号在当前通道已存在，不允许重复添加，请核实", zap.Any("ac_account", vca.AcAccount))
			return errors.New("账号在当前通道已存在，不允许重复添加，请核查")
		}
		vca.AcAccount = utils.Trim(vca.AcAccount)
	}

	err = global.GVA_DB.Create(vca).Error
	//vca传入的所有值 转化成 vcaDB vbox.ChannelAccount存放

	if vca.Status == 1 {
		go func() {
			var vcaDB vbox.ChannelAccount
			err = global.GVA_DB.Model(vbox.ChannelAccount{}).Where("id =?", vca.ID).First(&vcaDB).Error

			conn, err := mq.MQ.ConnPool.GetConnection()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)

			ch, err := conn.Channel()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
			}

			body := http2.DoGinContextBody(c)

			oc := vboxReq.ChanAccAndCtx{
				Obj: vcaDB,
				Ctx: vboxReq.Context{
					Body:      string(body),
					ClientIP:  c.ClientIP(),
					Method:    c.Request.Method,
					UrlPath:   c.Request.URL.Path,
					UserAgent: c.Request.UserAgent(),
					UserID:    int(vcaDB.CreatedBy),
				},
			}
			marshal, err := json.Marshal(oc)

			err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)
		}()
	}

	return err
}

// DeleteChannelAccount 删除通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) DeleteChannelAccount(vca vbox.ChannelAccount, c *gin.Context) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var vcaDB vbox.ChannelAccount
		if err := tx.Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).First(&vcaDB).Error; err != nil {
			return err
		}

		conn, err := mq.MQ.ConnPool.GetConnection()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)

		ch, err := conn.Channel()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		body := http2.DoGinContextBody(c)

		oc := vboxReq.ChanAccAndCtx{
			Obj: vcaDB,
			Ctx: vboxReq.Context{
				Body:      string(body),
				ClientIP:  c.ClientIP(),
				Method:    c.Request.Method,
				UrlPath:   c.Request.URL.Path,
				UserAgent: c.Request.UserAgent(),
				UserID:    int(vca.DeletedBy),
			},
		}
		marshal, err := json.Marshal(oc)

		err = ch.Publish(task.ChanAccDelCheckExchange, task.ChanAccDelCheckKey, marshal)

		if err := tx.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).Update("sys_status", 2).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

// DeleteChannelAccountByIds 批量删除通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) DeleteChannelAccountByIds(ids request.IdsReq, c *gin.Context, deletedBy uint) (err error) {

	if len(ids.Ids) < 1 {
		return fmt.Errorf("传入的id为空")
	} else {
		for _, ID := range ids.Ids {
			err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				var vcaDB vbox.ChannelAccount
				if err := tx.Model(&vbox.ChannelAccount{}).Where("id = ?", ID).First(&vcaDB).Error; err != nil {
					return err
				}

				conn, err := mq.MQ.ConnPool.GetConnection()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
				}
				defer mq.MQ.ConnPool.ReturnConnection(conn)

				ch, err := conn.Channel()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
				}

				body := http2.DoGinContextBody(c)

				oc := vboxReq.ChanAccAndCtx{
					Obj: vcaDB,
					Ctx: vboxReq.Context{
						Body:      string(body),
						ClientIP:  c.ClientIP(),
						Method:    c.Request.Method,
						UrlPath:   c.Request.URL.Path,
						UserAgent: c.Request.UserAgent(),
						UserID:    int(deletedBy),
					},
				}
				marshal, err := json.Marshal(oc)

				err = ch.Publish(task.ChanAccDelCheckExchange, task.ChanAccDelCheckKey, marshal)

				if err := tx.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("sys_status", 2).Error; err != nil {
					return err
				}

				return nil
			})
		}
	}

	return err
}

// SwitchEnableChannelAccount 开关通道账号
// Author [bboyhan](https://github.com/bboyhan)
func (vcaService *ChannelAccountService) SwitchEnableChannelAccount(vca vboxReq.ChannelAccountUpd, c *gin.Context) (err error) {
	var vcaDB vbox.ChannelAccount
	err = global.GVA_DB.Where("id = ?", vca.ID).First(&vcaDB).Error
	if err != nil {
		return fmt.Errorf("不存在的账号，请核查")
	}

	// 如果是开启，则发起一条消息，去查这个账号是否能开启

	//go func() {
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	ch, err := conn.Channel()
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
	}

	body := http2.DoGinContextBody(c)
	vcaDB.Status = vca.Status
	oc := vboxReq.ChanAccAndCtx{
		Obj: vcaDB,
		Ctx: vboxReq.Context{
			Body:      string(body),
			ClientIP:  c.ClientIP(),
			Method:    c.Request.Method,
			UrlPath:   c.Request.URL.Path,
			UserAgent: c.Request.UserAgent(),
			UserID:    int(vcaDB.CreatedBy),
		},
	}
	marshal, err := json.Marshal(oc)

	err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)
	//}()

	err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).Update("status", vca.Status).Update("updated_by", vca.UpdatedBy).Error
	return err
}

// SwitchEnableChannelAccountByIds 批量开关通道账号记录
// Author [bboyhan](https://github.com/bboyhan)
func (vcaService *ChannelAccountService) SwitchEnableChannelAccountByIds(upd vboxReq.ChannelAccountUpd, updatedBy uint, c *gin.Context) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

		// 如果是开启，则发起一条消息，去查这个账号是否能开启

		go func() {
			conn, err := mq.MQ.ConnPool.GetConnection()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)
			var vcaDBList []vbox.ChannelAccount
			err = global.GVA_DB.Model(vbox.ChannelAccount{}).Where("id in ?", upd.Ids).Find(&vcaDBList).Error

			for _, vcaDB := range vcaDBList {

				ch, err := conn.Channel()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
					continue
				}

				body := http2.DoGinContextBody(c)
				vcaDB.Status = upd.Status

				oc := vboxReq.ChanAccAndCtx{
					Obj: vcaDB,
					Ctx: vboxReq.Context{
						Body:      string(body),
						ClientIP:  c.ClientIP(),
						Method:    c.Request.Method,
						UrlPath:   c.Request.URL.Path,
						UserAgent: c.Request.UserAgent(),
						UserID:    int(vcaDB.CreatedBy),
					},
				}
				marshal, err := json.Marshal(oc)

				err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)
			}
		}()

		if err := tx.Unscoped().Model(&vbox.ChannelAccount{}).Where("id in ?", upd.Ids).Update("status", upd.Status).Update("updated_by", updatedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", upd.Ids).Updates(&vbox.ChannelAccount{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelAccount 更新通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) UpdateChannelAccount(vca vbox.ChannelAccount) (err error) {
	token := vca.Token
	//增加校验
	if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) || global.DnfContains(vca.Cid) {
		_, _, errX := product.Secret(token)
		if errX != nil {
			return errX
		}
	} else if global.SdoContains(vca.Cid) {
		parsedURL, errX := url.Parse(token)
		if errX != nil {
			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
			return errors.New("无效的URL")
		}
		query := parsedURL.Query()
		sndaId := query.Get("sndaId")
		if sndaId == "" {
			return errors.New("账号信息不完整")
		}
	} else if global.ECContains(vca.Cid) {
		b := http2.IsValidCookie(token)
		if !b {
			return errors.New("传入的ck不合法，请核查")
		}
	} else if global.QNContains(vca.Cid) {
		isCK := http2.IsValidCookie(token)
		if !isCK {
			return errors.New("ck信息不合法")
		}
		account, errX := product.FindAccNick(token)
		if errX != nil {
			global.GVA_LOG.Warn("未能解析到nick lid的值， 进行随机赋值")
		}
		vca.AcAccount = account

	} else if global.J3Contains(vca.Cid) {
		parsedURL, errX := url.Parse(token)
		if errX != nil {
			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
			return errors.New("无效的URL")
		}

		query := parsedURL.Query()
		account := query.Get("account")
		zoneCode := query.Get("zoneCode")
		SN := query.Get("SN")
		sign := query.Get("sign")
		if account == "" || zoneCode == "" || SN == "" || sign == "" {
			return errors.New("账号信息不完整")
		}
		if vca.Cid == "2001" && zoneCode != "z22" {
			return errors.New("仅支持双线二区参数，请核查")
		}
		if vca.Cid == "2002" && zoneCode != "z05" {
			return errors.New("仅支持电信五区参数，请核查")
		}
		vca.AcAccount = account
	} else if global.DyContains(vca.Cid) {
		b := http2.IsValidCookie(token)
		if !b {
			return errors.New("传入的ck不合法，请核查")
		}
	} else {
		return errors.New("该信道暂不支持更新账号")
	}
	err = global.GVA_DB.Save(&vca).Error
	return err
}

// GetChannelAccount 根据id获取通道账号记录
func (vcaService *ChannelAccountService) GetChannelAccount(id uint) (vca vbox.ChannelAccount, err error) {
	err = global.GVA_DB.Unscoped().Where("id = ?", id).First(&vca).Error
	var sysUser sysModel.SysUser
	err = global.GVA_DB.Unscoped().Where("id = ?", vca.CreatedBy).First(&sysUser).Error
	vca.Username = sysUser.Username

	rdAccId := fmt.Sprintf(global.AccRecord, vca.AcId)
	var records []sysModel.SysOperationRecord
	global.GVA_DB.Model(&sysModel.SysOperationRecord{}).Distinct("resp,created_at").Where("mark_id = ?", rdAccId).Scan(&records)

	ext := map[string]interface{}{
		"records": records,
	}
	vca.Ext = ext
	return
}

// GetChannelAccountByAcId 根据AcId获取通道账号记录
func (vcaService *ChannelAccountService) GetChannelAccountByAcId(acId string) (vca vbox.ChannelAccount, err error) {
	err = global.GVA_DB.Unscoped().Where("ac_id = ?", acId).First(&vca).Error
	var sysUser sysModel.SysUser
	err = global.GVA_DB.Unscoped().Where("id = ?", vca.CreatedBy).First(&sysUser).Error
	vca.Username = sysUser.Username

	rdAccId := fmt.Sprintf(global.AccRecord, acId)
	var records []sysModel.SysOperationRecord
	global.GVA_DB.Model(&sysModel.SysOperationRecord{}).Distinct("resp,created_at").Where("mark_id = ?", rdAccId).Scan(&records)

	ext := map[string]interface{}{
		"records": records,
	}
	vca.Ext = ext
	return
}

// GetChannelAccountInfoList 分页获取通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) GetChannelAccountInfoList(info vboxReq.ChannelAccountSearch, ids []uint) (list []vbox.ChannelAccount, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelAccount{})
	var vcas []vbox.ChannelAccount
	db.Where("created_by in (?)", ids)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Username != "" {
		var sysUserIDs []uint
		err = global.GVA_DB.Unscoped().Table("sys_users").Select("id").Where("username LIKE ?", "%"+info.Username+"%").Scan(&sysUserIDs).Error
		if len(sysUserIDs) > 0 {
			db = db.Where("created_by in (?)", sysUserIDs)
		}
	}
	if info.AcRemark != "" {
		db = db.Where("ac_remark LIKE ?", info.AcRemark+"%")
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account LIKE ?", info.AcAccount+"%")
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysStatus != nil {
		db = db.Where("sys_status = ?", info.SysStatus)
	}
	if info.CtlStatus != nil {
		db = db.Where("ctl_status = ?", info.CtlStatus)
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("id desc").Find(&vcas).Error
	return vcas, total, err
}
