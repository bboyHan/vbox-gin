package task

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func HandleAccAvailable() (err error) {
	//自定义 通道账号可用检测
	rdConn := global.GVA_REDIS.Conn()
	defer rdConn.Close()
	var idList []uint
	// 拿出现在所有付方可用的账户
	err = global.GVA_DB.Model(&vbox.PayAccount{}).Table("vbox_pay_account").
		Select("uid").Where("status = ?", 1).Find(&idList).Error
	if err != nil {
		global.GVA_LOG.Error("查付方库数据异常", zap.Error(err))
		return
	}

	//global.GVA_LOG.Info("根据开启的商户列表，开始检测可用账号情况")

	for _, uid := range idList {
		var channelCodeList []string

		// 获取组织ID
		orgTmp := utils2.GetSelfOrg(uid)
		orgID := orgTmp[0]
		// 当前付方所拥有的产品列表
		chanKey := fmt.Sprintf(global.OrgChanSet, orgID)

		c, err := rdConn.Exists(context.Background(), chanKey).Result()
		if c == 0 {
			var productIds []uint
			if err != nil {
				global.GVA_LOG.Error("当前缓存池无此用户对应的orgIds，redis err", zap.Error(err))
			}
			if err = global.GVA_DB.Model(&vbox.OrgProduct{}).Distinct("channel_product_id").Select("channel_product_id").Where("organization_id = ?", orgID).Find(&productIds).Error; err != nil {
				global.GVA_LOG.Error("OrgProduct查该组织下数据channel code异常", zap.Error(err))
				continue
			}
			if err = global.GVA_DB.Model(&vbox.ChannelProduct{}).Select("channel_code").Where("id in ?", productIds).Find(&channelCodeList).Error; err != nil {
				global.GVA_LOG.Error("ChannelProduct查channelCodeList 库数据异常", zap.Error(err))
				continue
			}

			//jsonStr, _ := json.Marshal(channelCodeList)
			for _, cid := range channelCodeList {
				rdConn.SAdd(context.Background(), chanKey, cid)
			}
		} else {
			cidList, _ := rdConn.SMembers(context.Background(), chanKey).Result()
			//err = json.Unmarshal(jsonStr, &channelCodeList)
			channelCodeList = cidList
		}

		// 遍历每个归属产品下的通道账号情况，筛选出可用的账号 入等待池
		for _, channelCode := range channelCodeList {
			deepUIDs := utils2.GetSelfOrg(uid)
			var moneyList []string
			var vcas []vbox.ChannelAccount

			err := global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
				Where("created_by in ?", deepUIDs).Where("cid = ?", channelCode).Where("status = ? and sys_status = ?", 1, 1).Scan(&vcas).Error

			if err != nil {
				global.GVA_LOG.Error("ChannelAccount查数据异常", zap.Error(err))
				continue
			}
			if len(vcas) == 0 {
				//global.GVA_LOG.Warn("ChannelAccount查数据，当前产品编码无开启的账号", zap.Any("channel code", channelCode))
				continue
			} else {
				//	遍历可用的账号，查一下官方记录和库里的订单情况
				duration, _ := HandleExpTime2Product(channelCode)
				nowTime := time.Now()

				// 当前时间前后各缓冲1分钟时间
				startTime := nowTime.Add(-duration).Add(-60 * time.Second)
				endTime := nowTime.Add(60 * time.Second)
				//global.GVA_LOG.Info("查询时间范围", zap.Any("startTime", startTime), zap.Any("endTime", endTime))

				moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgID, channelCode)
				cm, err := rdConn.Exists(context.Background(), moneyKey).Result()
				if cm == 0 {
					if err != nil {
						global.GVA_LOG.Error("redis err", zap.Error(err))
					}

					// 0 - 如果是引导的，金额从开启的商铺金额获取
					// 1 - 如果是预产的，从未过期的预产金额获取

					// 0 -
					if global.TxContains(channelCode) { // 10xx tx引导
						err := global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
							Where("cid = ? and status = ? and created_by in ?", channelCode, 1, deepUIDs).
							Find(&moneyList).Error
						if err != nil {
							global.GVA_LOG.Error("查库ex", zap.Error(err))
						} else {
							//jsonStr, _ := json.Marshal(moneyList)
							for _, money := range moneyList {
								rdConn.SAdd(context.Background(), moneyKey, money)
							}
							//rdConn.Set(context.Background(), moneyKey, jsonStr, 10*time.Minute)
						}
					} else if global.J3Contains(channelCode) { // 10xx 剑三引导
						err := global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
							Where("cid = ? and status = ? and created_by in ?", channelCode, 1, deepUIDs).
							Find(&moneyList).Error
						if err != nil {
							global.GVA_LOG.Error("查库ex", zap.Error(err))
						} else {
							for _, money := range moneyList {
								rdConn.SAdd(context.Background(), moneyKey, money)
							}
							//jsonStr, _ := json.Marshal(moneyList)
							//rdConn.SAdd(context.Background(), moneyKey, jsonStr, 10*time.Minute)
						}
						//	1 -
					} else if global.PcContains(channelCode) { // 30xx 付码
						err := global.GVA_DB.Model(&vbox.ChannelPayCode{}).Distinct("money").Select("money").
							Where("cid = ? and code_status = ? and created_by in ?", channelCode, 2, deepUIDs).
							Find(&moneyList).Error
						if err != nil {
							global.GVA_LOG.Error("查库ex", zap.Error(err))
						} else {
							for _, money := range moneyList {
								rdConn.SAdd(context.Background(), moneyKey, money)
							}
						}
					}

				} else {
					moneyMembers, _ := rdConn.SMembers(context.Background(), moneyKey).Result()
					moneyList = moneyMembers
				}

				//global.GVA_LOG.Info("ddd - ", zap.Any("cid", channelCode), zap.Any("moneyList", moneyList))

				for _, vca := range vcas {
					vcaTmp := vca
					cid := channelCode
					go func() {
						accMem := strconv.FormatUint(uint64(vcaTmp.ID), 10) + "_" + vcaTmp.AcId + "_" + vcaTmp.AcAccount
						if global.TxContains(channelCode) {

							//	查tx 记录
							records, err := product.QryQQRecordsBetween(vcaTmp, startTime, endTime)
							if err != nil {
								// 查单有问题，直接订单要置为超时，消息置为处理完毕
								global.GVA_LOG.Error("查单异常跳过", zap.Error(err))
								return
							}
							rdMap := product.Classifier(records.WaterList)

							if vm, ok := rdMap["Q币"]; !ok {
								//global.GVA_LOG.Info("还没有QB的充值记录")
								for _, money := range moneyList {
									key := fmt.Sprintf(global.ChanOrgAccZSet, orgID, cid, money)

									// 再查一下库，这个时间段的有没有这个账号的订单
									var count int64
									err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("created_at between ? and ?", startTime, endTime).
										Where("ac_id = ? and money = ?", vcaTmp.AcId, money).
										Count(&count).Error
									if err != nil {
										global.GVA_LOG.Error("查库异常", zap.Error(err))
									}
									if count > 0 {
										//global.GVA_LOG.Info(fmt.Sprintf("库里有这个时间段的订单数据, count: [%d]", count))
										// 证明这种金额的，已经有记录了，代表这种金额对应的账号暂不可用
										global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
											Score:  1,
											Member: accMem,
										})
									} else {

										// 进入等待拿走的可用账号池子
										global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
											Score:  0,
											Member: accMem,
										})
									}
								}
							} else { // 有qb 记录了，要查一下qb所充的金额有没有相对应的存在

								for _, money := range moneyList {
									key := fmt.Sprintf(global.ChanOrgAccZSet, orgID, cid, money)
									if rd, ok2 := vm[money+"00"]; !ok2 {
										global.GVA_LOG.Info(fmt.Sprintf("还没有QB的充值记录, ac account : [%s], 金额: [%s]", rd, money))

										// 再查一下库，这个时间段的有没有这个账号的订单
										var count int64
										err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("created_at between ? and ?", startTime, endTime).
											Where("ac_id = ? and money = ?", vcaTmp.AcId, money).
											Count(&count).Error
										if err != nil {
											global.GVA_LOG.Error("查库异常")
										}
										if count > 0 {
											global.GVA_LOG.Info(fmt.Sprintf("库里有这个时间段的订单数据, count: [%d]", count))
											// 证明这种金额的，已经有记录了，代表这种金额对应的账号暂不可用
											global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
												Score:  1,
												Member: accMem,
											})
										} else { // 查了库也没数据，就可以加到可用列表中待取用了
											// 进入等待拿走的可用账号池子
											global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
												Score:  0,
												Member: accMem,
											})
										}

									} else { // 证明这种金额的，已经有记录了，代表这种金额对应的账号暂不可用
										global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
											Score:  1,
											Member: accMem,
										})
									}
								}

							}

						} else if global.J3Contains(channelCode) {
							// 查j3 记录
						} else if global.PcContains(channelCode) {
							//	查付码记录

							//	查tx 记录
							records, err := product.QryQQRecordsBetween(vcaTmp, startTime, endTime)
							if err != nil {
								// 查单有问题，直接订单要置为超时，消息置为处理完毕
								global.GVA_LOG.Error("查单异常跳过", zap.Error(err))
								return
							}
							rdMap := product.Classifier(records.WaterList)

							if vm, ok := rdMap["Q币"]; !ok {
								//global.GVA_LOG.Info("还没有QB的充值记录")
								for _, money := range moneyList {
									key := fmt.Sprintf(global.ChanOrgAccZSet, orgID, cid, money)
									//global.GVA_LOG.Info("打印出来了", zap.Any("key", key))

									// 再查一下库，这个时间段的有没有这个账号的订单
									var count int64
									err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("created_at between ? and ?", startTime, endTime).
										Where("ac_id = ? and money = ?", vcaTmp.AcId, money).
										Count(&count).Error
									if err != nil {
										global.GVA_LOG.Error("查库异常", zap.Error(err))
									}
									if count > 0 {
										// 大于0证明已经有充值数据了，置为score 1，代表不可用
										global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
											Score:  1,
											Member: accMem,
										})

									} else {

										// 进入等待拿走的可用账号池子
										global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
											Score:  0,
											Member: accMem,
										})
									}
								}
							} else { // 有qb 记录了，要查一下qb所充的金额有没有相对应的存在

								for _, money := range moneyList {
									key := fmt.Sprintf(global.ChanOrgAccZSet, orgID, cid, money)
									if rd, ok2 := vm[money+"00"]; !ok2 {
										global.GVA_LOG.Info(fmt.Sprintf("还没有QB的充值记录, ac account : [%s], 金额: [%s]", rd, money))

										// 再查一下库，这个时间段的有没有这个账号的订单
										var count int64
										err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("created_at between ? and ?", startTime, endTime).
											Where("ac_id = ? and money = ?", vcaTmp.AcId, money).
											Count(&count).Error
										if err != nil {
											global.GVA_LOG.Error("查库异常")
										}
										if count > 0 {
											global.GVA_LOG.Info(fmt.Sprintf("库里有这个时间段的订单数据, count: [%d]", count))
											// 大于0证明已经有充值数据了，置为score 1，代表不可用
											global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
												Score:  1,
												Member: accMem,
											})

										} else { // 查了库也没数据，就可以加到可用列表中待取用了
											// 进入等待拿走的可用账号池子
											global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
												Score:  0,
												Member: accMem,
											})
										}

									} else { // 证明这种金额的，已经有记录了，代表这种金额对应的账号暂不可用
										global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
											Score:  1,
											Member: accMem,
										})
									}
								}

							}
						}
					}()

				}

			}

			//global.GVA_LOG.Info("当前可用账号情况", zap.String("channel code", channelCode), zap.Any("可用数", len(vcas)), zap.Any("list", vcas))
		}
	}

	return err
}

func HandleExpTime2Product(chanID string) (time.Duration, error) {
	var key string

	if global.TxContains(chanID) {
		key = "1000"
	} else if global.J3Contains(chanID) {
		key = "2000"
	} else if global.PcContains(chanID) {
		key = "3000"
	}

	var expTimeStr string
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", key).
			First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return 0, err
		}
		expTimeStr = proxy.Url
		seconds, _ := strconv.Atoi(expTimeStr)
		duration := time.Duration(seconds) * time.Second

		global.GVA_REDIS.Set(context.Background(), key, int64(duration.Seconds()), 0)
		global.GVA_LOG.Info("数据库取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))

		return duration, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
		return 0, err
	} else {
		expTimeStr, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		seconds, _ := strconv.Atoi(expTimeStr)

		duration := time.Duration(seconds) * time.Second

		//global.GVA_LOG.Info("缓存池取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))
		return duration, err
	}
}
