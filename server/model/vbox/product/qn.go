package product

// prePageNo=1&sifg=0&action=itemlist%2FSoldQueryAction&tabCode=success&rateStatus=&orderStatus=SUCCESS&payDateBegin=0&buyerNick=&pageSize=15&dateEnd=0&rxOldFlag=0&rxSendFlag=0&useCheckcode=false&dateBegin=0&tradeTag=0&rxHasSendFlag=0&auctionType=0&close=0&sellerNick=&notifySendGoodsType=ALL&sellerMemoFlag=0&useOrderInfo=false&logisticsService=&isQnNew=true&pageNum=1&o2oDeliveryType=ALL&rxAuditFlag=0&queryOrder=desc&rxElectronicAuditFlag=0&queryMore=false&payDateEnd=0&rxWaitSendflag=0&sellerMemo=0&rxElectronicAllFlag=0&rxSuccessflag=0&refund=&errorCheckcode=false&isHideNick=true
type QueryData struct {
	PrePageNo             int    `json:"prePageNo"`
	Sifg                  int    `json:"sifg"`
	Action                string `json:"action"`
	TabCode               string `json:"tabCode"`
	RateStatus            string `json:"rateStatus"`
	OrderStatus           string `json:"orderStatus"`
	PayDateBegin          int    `json:"payDateBegin"`
	BuyerNick             string `json:"buyerNick"`
	PageSize              int    `json:"pageSize"`
	DateEnd               int    `json:"dateEnd"`
	RxOldFlag             int    `json:"rxOldFlag"`
	RxSendFlag            int    `json:"rxSendFlag"`
	UseCheckcode          bool   `json:"useCheckcode"`
	DateBegin             int    `json:"dateBegin"`
	TradeTag              int    `json:"tradeTag"`
	RxHasSendFlag         int    `json:"rxHasSendFlag"`
	AuctionType           int    `json:"auctionType"`
	Close                 int    `json:"close"`
	SellerNick            string `json:"sellerNick"`
	NotifySendGoodsType   string `json:"notifySendGoodsType"`
	SellerMemoFlag        int    `json:"sellerMemoFlag"`
	UseOrderInfo          bool   `json:"useOrderInfo"`
	LogisticsService      string `json:"logisticsService"`
	IsQnNew               bool   `json:"isQnNew"`
	PageNum               int    `json:"pageNum"`
	O2oDeliveryType       string `json:"o2oDeliveryType"`
	RxAuditFlag           int    `json:"rxAuditFlag"`
	QueryOrder            string `json:"queryOrder"`
	RxElectronicAuditFlag int    `json:"rxElectronicAuditFlag"`
	QueryMore             bool   `json:"queryMore"`
	PayDateEnd            int    `json:"payDateEnd"`
	RxWaitSendflag        int    `json:"rxWaitSendflag"`
	SellerMemo            int    `json:"sellerMemo"`
	RxElectronicAllFlag   int    `json:"rxElectronicAllFlag"`
	RxSuccessflag         int    `json:"rxSuccessflag"`
	Refund                string `json:"refund"`
	ErrorCheckcode        bool   `json:"errorCheckcode"`
	IsHideNick            bool   `json:"isHideNick"`
}

type QnResp struct {
	TraceId         string `json:"traceId"`
	IsQnNewHideNick bool   `json:"isQnNewHideNick"`
	Extra           struct {
		BatchFlagUrl             int    `json:"batchFlagUrl"`
		SellerFlag               int    `json:"sellerFlag"`
		CurrencySymbol           string `json:"currencySymbol"`
		IsQnNew                  bool   `json:"isQnNew"`
		IsShowSellerService      bool   `json:"isShowSellerService"`
		IsQnNewExportNew         bool   `json:"isQnNewExportNew"`
		IsQnNewHideNick          bool   `json:"isQnNewHideNick"`
		IsQnNewShowSellerService bool   `json:"isQnNewShowSellerService"`
	} `json:"extra"`
}

type QnOrderRecord struct {
	//"statusInfo": {
	//                "operations": [
	//                    {
	//                        "style": "t16",
	//                        "text": "详情",
	//                        "params": {
	//                            "qianNiuPCDetailUrl": "//qn.taobao.com/home.htm/trade-platform/tp/detail?bizOrderId=3668803670628269143",
	//                            "solutionPCDetailUrl": "//trade.taobao.com/trade/detail/trade_item_detail.htm?bizOrderId=3668803670628269143"
	//                        },
	//                        "type": "operation",
	//                        "url": "//trade.taobao.com/trade/detail/trade_item_detail.htm?bizOrderId=3668803670628269143"
	//                    },
	//                    {
	//                        "style": "t16",
	//                        "id": "printWaybill",
	//                        "text": "打单",
	//                        "type": "operation",
	//                        "url": "https://myseller.taobao.com/home.htm/qn-order/unshipped?ORDER_ID=3668803670628269143"
	//                    },
	//                    {
	//                        "style": "t15",
	//                        "text": "补发货",
	//                        "type": "operation",
	//                        "url": "//qn.taobao.com/home.htm/consign-order/?from=list&tradeId=3668803670628269143&type=1"
	//                    }
	//                ],
	//                "text": "交易成功",
	//                "type": "t0"
	//            }
	StatusInfo struct {
		Operations []struct {
			ID     string `json:"id"`
			Style  string `json:"style"`
			Text   string `json:"text"`
			Params struct {
				QianNiuPCDetailUrl  string `json:"qianNiuPCDetailUrl"`
				SolutionPCDetailUrl string `json:"solutionPCDetailUrl"`
			} `json:"params"`
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"operations"`
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"statusInfo"`
	//"extra": {
	//                "sellerFlag": 0,
	//                "currencySymbol": "￥",
	//                "inHold": false,
	//                "isShowSellerService": false,
	//                "currency": "CNY",
	//                "sellerFlagPic": "https://gw.alicdn.com/imgextra/i3/O1CN01kOL67z1j0BDjYiCKU_!!6000000004485-2-tps-120-120.png?getAvatar=avatar",
	//                "batchSendGoods": 1,
	//                "disableCheckbox": false
	//            },
	Extra struct {
		SellerFlag          int    `json:"sellerFlag"`
		CurrencySymbol      string `json:"currencySymbol"`
		InHold              bool   `json:"inHold"`
		IsShowSellerService bool   `json:"isShowSellerService"`
		Currency            string `json:"currency"`
		SellerFlagPic       string `json:"sellerFlagPic"`
		BatchSendGoods      int    `json:"batchSendGoods"`
		DisableCheckbox     bool   `json:"disableCheckbox"`
	}

	OrderInfo struct {
		CreateTime string `json:"createTime"`
		Id         string `json:"id"`
	} `json:"orderInfo"`
	Id string `json:"id"`
}
