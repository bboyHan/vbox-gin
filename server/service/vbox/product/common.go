package product

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	"reflect"
	"strconv"
	"strings"
)

// Classifier 计算不同类型 - 不同金额 - 记录集合
func Classifier(payments interface{}) map[string]map[string][]string {
	// 使用map存储不同充值类型下的支付金额和充值账号ID集合（去重）
	paymentsByTypeAndAmount := make(map[string]map[string][]string)

	switch v := payments.(type) {
	case []product.Payment:
		for _, payment := range payments.([]product.Payment) {
			amount := payment.PayAmt
			showName := payment.ShowName
			provideID := payment.ProvideID
			if strings.Contains(showName, "DNF") {
				showName = "DNF"
			}

			// 检查是否存在对应的充值类型的map
			if _, ok := paymentsByTypeAndAmount[showName]; !ok {
				paymentsByTypeAndAmount[showName] = make(map[string][]string)
			}

			// 添加充值账号ID到对应的支付金额中（去重）
			ids := paymentsByTypeAndAmount[showName][amount]
			exists := false
			for _, id := range ids {
				if id == provideID {
					exists = true
					break
				}
			}
			if !exists {
				paymentsByTypeAndAmount[showName][amount] = append(ids, provideID)
			}
		}
	case []product.SdoDaoYuOrderRecord:
		for _, payment := range payments.([]product.SdoDaoYuOrderRecord) {
			orderAmount := strconv.FormatFloat(payment.OrderAmount, 'f', -1, 64)
			parts := strings.Split(orderAmount, ".")

			var intStrAmount string
			if len(parts) > 0 {
				intStrAmount = parts[0]
			} else {
				intStrAmount = orderAmount
			}

			appName := payment.AppName
			accAndOrderID := fmt.Sprintf("%s_%s", payment.DisplayAccount, payment.OrderId)
			if payment.PayStatus != 5 {
				continue
			}

			// 检查是否存在对应的充值类型的map
			if _, ok := paymentsByTypeAndAmount[appName]; !ok {
				paymentsByTypeAndAmount[appName] = make(map[string][]string)
			}

			// 添加充值账号ID到对应的支付金额中（去重）
			ids := paymentsByTypeAndAmount[appName][intStrAmount]
			exists := false
			for _, id := range ids {
				if id == accAndOrderID {
					exists = true
					break
				}
			}
			if !exists {
				paymentsByTypeAndAmount[appName][intStrAmount] = append(ids, accAndOrderID)
			}
		}
	case []product.SdoOrderRecord:
		for _, payment := range payments.([]product.SdoOrderRecord) {
			orderAmount := payment.OrderAmount
			parts := strings.Split(orderAmount, ".")

			var intStrAmount string
			if len(parts) > 0 {
				intStrAmount = parts[0]
			} else {
				intStrAmount = orderAmount
			}

			appName := payment.AppName
			accAndOrderID := fmt.Sprintf("%s_%s", payment.InputOrderUser, payment.OrderID)
			if payment.State != 5 {
				continue
			}

			// 检查是否存在对应的充值类型的map
			if _, ok := paymentsByTypeAndAmount[appName]; !ok {
				paymentsByTypeAndAmount[appName] = make(map[string][]string)
			}

			// 添加充值账号ID到对应的支付金额中（去重）
			ids := paymentsByTypeAndAmount[appName][intStrAmount]
			exists := false
			for _, id := range ids {
				if id == accAndOrderID {
					exists = true
					break
				}
			}
			if !exists {
				paymentsByTypeAndAmount[appName][intStrAmount] = append(ids, accAndOrderID)
			}
		}
	default:
		// Handle other types or provide an error message
		fmt.Println("Unsupported type:", reflect.TypeOf(v))
	}

	// 输出结果
	//for showName, amounts := range paymentsByTypeAndAmount {
	//	fmt.Printf("充值类型：%s\n", showName)
	//	for amount, ids := range amounts {
	//		fmt.Printf("支付金额：%s，充值账号ID集合：%v\n", amount, ids)
	//	}
	//	fmt.Println()
	//}
	return paymentsByTypeAndAmount
}
