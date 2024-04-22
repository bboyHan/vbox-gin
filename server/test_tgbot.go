package main

//import (
//	"encoding/json"
//	"fmt"
//	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//	"log"
//	"strconv"
//	"time"
//
//	"gopkg.in/telebot.v3"
//)
//
//var headers = map[string]string{
//	"Accept-Encoding": "deflate",
//	"Accept-Language": "zh-CN,zh;q=0.9",
//	"Connection":      "keep-alive",
//	//Cookie: sl-session=61W9AMrdIGY8hksHYlzN1g==; __51cke__=; __tins__21838175=%7B%22sid%22%3A%201713343564248%2C%20%22vd%22%3A%202%2C%20%22expires%22%3A%201713345375293%7D; __51laig__=2
//	"Host":             "www.filgd.com",
//	"Referer":          "http://www.filgd.com/jnd28.html",
//	"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
//	"X-Requested-With": "XMLHttpRequest",
//}
//var options = &vbHttp.RequestOptions{
//	Headers:      headers,
//	MaxRedirects: 3,
//	PayloadType:  "url",
//}
//
///*
//	{
//	    "type": "jnd28",
//	    "sf": "1",
//	    "ms": "zh",
//	    "issue": [
//	        {
//	            "qishu": "3129086",
//	            "time": "04-17 18:21:30",
//	            "num1": "3",
//	            "num2": "9",
//	            "num3": "1",
//	            "sum": "13",
//	            "dx": "小",
//	            "ds": "单",
//	            "next": 1713349500
//	        }
//		"time": {
//			"t": 195,
//			"h": "00",
//			"m": "03",
//			"s": 15
//		},
//		"nowTime": "2024-04-18"
//
//]
//*/
//type R struct {
//	Type  string `json:"type"`
//	Sf    string `json:"sf"`
//	Ms    string `json:"ms"`
//	Issue []struct {
//		Qishu string `json:"qishu"`
//		Time  string `json:"time"`
//		Num1  string `json:"num1"`
//		Num2  string `json:"num2"`
//		Num3  string `json:"num3"`
//		Sum   string `json:"sum"`
//		Dx    string `json:"dx"`
//		Ds    string `json:"ds"`
//		Next  int64  `json:"next"`
//	} `json:"issue"`
//	Time struct {
//		T int    `json:"t"`
//		H string `json:"h"`
//		M string `json:"m"`
//		S int    `json:"s"`
//	} `json:"time"`
//	NowTime string `json:"nowTime"`
//}
//
//var (
//	r = &telebot.ReplyMarkup{}
//
//	balanceBtn = r.Text("余额查询")
//	groupBtn   = r.URL("上下分群", "https://t.me/pc666_bot")
//	hisBtn     = r.URL("投注记录", "https://t.me/pc666_bot")
//	incomeBtn  = r.URL("自助充值", "https://t.me/pc666_bot")
//)
//
//func main() {
//	pref := telebot.Settings{
//		Token:  "6858573723:AAGTlf__LLngogpRfAQ5HBVa4C4mWlcexrc",
//		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
//	}
//
//	b, err := telebot.NewBot(pref)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//	b.Handle("hello", func(c telebot.Context) error {
//		client := vbHttp.NewHTTPClient()
//		//http://www.filgd.com/data/get/checkData?type=jnd28&sf=1&ms=zh
//		resp, _ := client.Get("http://www.filgd.com/data/get/checkData?type=jnd28&sf=1&ms=zh", options)
//		var rt R
//		err = json.Unmarshal(resp.Body, &rt)
//		var rs string
//		if err != nil {
//			log.Println(err)
//			return nil
//		}
//		if rt.Time.T < 10 {
//			return c.Send("===封盘中，下注无效===")
//		}
//		if len(rt.Issue) < 1 {
//			return nil
//		}
//		rt.Issue = rt.Issue[0:10]
//		f := rt.Issue[0]
//		qs, _ := strconv.ParseInt(f.Qishu, 10, 64)
//		newQS := qs + 1
//		waitQS := fmt.Sprintf("%d", newQS)
//
//		//开奖时间
//		rs += fmt.Sprintf("💎 `%s` 期 开始投注💎 \n 【⏰倒计时： %v 时 %v 分 %v 秒 】 \n -—— - ——- - -—— - ——- \n \n\n -—— - ——- 历史开奖 -—— - ——- \n", waitQS, rt.Time.H, rt.Time.M, rt.Time.S)
//
//		for _, v := range rt.Issue {
//			qishu := v.Qishu
//			num1 := v.Num1
//			num2 := v.Num2
//			num3 := v.Num3
//			//转为int加和
//			n1, _ := strconv.ParseInt(num1, 10, 64)
//			n2, _ := strconv.ParseInt(num2, 10, 64)
//			n3, _ := strconv.ParseInt(num3, 10, 64)
//			n := n1 + n2 + n3
//
//			rs += fmt.Sprintf("`%s`期 ：`%s + %s + %s = %d %s %s` \n", qishu, num1, num2, num3, n, v.Dx, v.Ds)
//		}
//
//		r.Inline(
//			r.Row(balanceBtn, groupBtn),
//			r.Row(hisBtn, incomeBtn),
//		)
//
//		opt := &telebot.SendOptions{ReplyMarkup: r, ParseMode: telebot.ModeMarkdown}
//
//		return c.Send(rs, opt)
//	})
//
//	b.Handle("/start", func(c telebot.Context) error {
//		return c.Send("Welcome to the bot!")
//	})
//
//	b.Handle(telebot.OnText, func(c telebot.Context) error {
//		// All the text messages that weren't
//		// captured by existing handlers.
//
//		var (
//			//user = c.Sender()
//			cc = c.Chat()
//			//text = c.Text()
//		)
//		p := &telebot.Photo{File: telebot.FromDisk("2.png")}
//
//		// Use full-fledged bot's functions
//		// only if you need a result:
//		_, err := b.Send(cc, p)
//		if err != nil {
//			return err
//		}
//
//		// Instead, prefer a context short-hand:
//		return nil
//
//		// Instead, prefer a context short-hand:
//		//return c.SendAlbum(to, Album{photo, &photo2}, ModeHTML)
//	})
//
//	b.Handle(telebot.OnText, func(c telebot.Context) error {
//		// All the text messages that weren't
//		// captured by existing handlers.
//
//		var (
//			user = c.Sender()
//			text = c.Text()
//		)
//
//		if text == "1押10" {
//			opt := &telebot.SendOptions{ReplyMarkup: r, ParseMode: telebot.ModeMarkdown}
//
//			_, err := b.Send(c.Chat(), user.Username+"下注: `1押10`", opt)
//			if err != nil {
//				return err
//			}
//			return nil
//		}
//
//		// Use full-fledged bot's functions
//		// only if you need a result:
//		_, err := b.Send(user, text)
//		if err != nil {
//			return err
//		}
//
//		// Instead, prefer a context short-hand:
//		return c.Send(user, text)
//	})
//
//	b.Start()
//}
