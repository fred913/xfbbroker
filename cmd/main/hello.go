package main

import (
	"fmt"
	"runtime/debug"

	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yiffyi/gorad"
	"github.com/yiffyi/gorad/notification"
	"github.com/yiffyi/xfbbroker"
	"github.com/yiffyi/xfbbroker/xfb"
)

var cfg *xfbbroker.Config

func rechargeToThreshold(curBalance float64, u *xfbbroker.User) error {
	if u.Threshold-curBalance >= 10 {
		delta := u.Threshold - curBalance
		if delta > 100 {
			delta = 100.0
		}

		payUrl, err := xfb.RechargeOnCard(strconv.FormatFloat(delta, 'f', 2, 64), u.OpenId, u.SessionId, u.YmUserId)
		if err != nil {
			slog.Error("unable to recharge", "err", err)
			return err
		}
		url, _ := url.Parse(payUrl)
		tranNo := url.Query().Get("tran_no")

		_, err = xfb.SignPayCheck(tranNo)
		if err != nil {
			slog.Error("signpay check failed", "err", err)
			return err
		}
		err = xfb.PayChoose(tranNo)
		if err != nil {
			slog.Error("choose signpay failed", "err", err)
			return err
		}
		err = xfb.DoPay(tranNo)
		if err != nil {
			slog.Error("unable to pay", "err", err)
			return err
		}
		slog.Info("recharge to balance", "name", u.Name, "delta", delta, "tranNo", tranNo)
	}
	return nil
}

func formatExpense(x string) string {
	f, err := strconv.ParseFloat(x, 64)
	if err != nil {
		slog.Error("formatExpense", "err", err)
		return "数据错误"
	}

	if f < 0 {
		return fmt.Sprintf("￥%.2f", -f)
	} else {
		return fmt.Sprintf("+￥%.2f", f)
	}
}

func needNotify(feeName string) bool {
	if feeName == "金额写卡" {
		return false
	}

	return true
}

func sendNotify(key string, t *xfb.Trans) error {
	if len(key) == 0 {
		return nil
	}
	bot := notification.WeComBot{Key: key}
	msg := map[string]interface{}{
		"msgtype": "template_card",
		"template_card": map[string]interface{}{
			"card_type": "text_notice",
			"source": map[string]interface{}{
				"desc": "校园卡账单",
			},
			"main_title": map[string]interface{}{
				"title": t.Address,
				"desc":  t.FeeName,
			},
			"emphasis_content": map[string]interface{}{
				"title": formatExpense(t.Money),
			},
			"horizontal_content_list": []map[string]string{
				{
					"keyname": "余额",
					"value":   t.AfterMon,
				},
				{
					"keyname": "流水号",
					"value":   t.Serialno,
				},
				{
					"keyname": "交易时间",
					"value":   t.Dealtime,
				},
				{
					"keyname": "到账时间",
					"value":   t.Time,
				},
			},
			"card_action": map[string]interface{}{
				"type": 1,
				"url":  cfg.AuthLocalUrl,
			},
		},
	}
	return bot.SendMessage(msg)
}

func sendError(key string, err error, u *xfbbroker.User) error {
	if len(key) == 0 {
		return nil
	}
	bot := notification.WeComBot{Key: key}
	msg := map[string]interface{}{
		"msgtype": "template_card",
		"template_card": map[string]interface{}{
			"card_type": "text_notice",
			"source": map[string]interface{}{
				"desc": "校园卡账单",
			},
			"main_title": map[string]interface{}{
				"title": "请求错误",
				"desc":  u.Name,
			},
			"sub_title_text": "自动轮询已取消，点击重新授权\n" + err.Error(),
			"horizontal_content_list": []map[string]string{
				{
					"keyname": "ymId",
					"value":   u.YmUserId,
				},
			},
			"card_action": map[string]interface{}{
				"type": 1,
				"url":  cfg.AuthLocalUrl,
			},
		},
	}
	return bot.SendMessage(msg)
}

func checkTransLoop() {
	ticker := time.NewTicker(time.Duration(cfg.CheckTransInterval) * time.Second)
	for {
		// select {
		// case <-ticker.C:
		needSave := false
		for k := range cfg.Users {
			u, _ := cfg.GetUser(k) // ensure locking
			if u.Enabled && u.Failed < 3 {
				total, rows, err := xfb.CardQuerynoPage(u.SessionId, u.YmUserId, time.Now())
				if err != nil {
					slog.Error("CardQuerynoPage failed", "err", err)
					goto fail
				} else {
					updated := false
					slog.Debug("check trans", "name", u.Name, "total", total)

					for i := len(rows) - 1; i >= 0; i-- {
						v := rows[i]
						s, err := strconv.Atoi(v.Serialno)
						if err != nil {
							slog.Error("bad Serialno", "err", err, "name", u.Name, "serial", v.Serialno)
							continue
						}
						if s > u.LastSerial {
							slog.Info("New transaction", "detail", v)

							if needNotify(v.FeeName) {
								err = sendNotify(u.WeComBotKey, &v)
								if err != nil {
									slog.Error("failed to notify", "err", err)
									break
								}
							} else {
								slog.Info("skipped", "feeName", v.FeeName)
							}

							if u.LastSerial < s {
								u.LastSerial = s
							}
							updated = true
						} else {
							continue
						}
					}

					if updated {
						goto set
					}
				}

				// success?
				continue
			fail:
				u.Failed++
				if u.Failed <= 3 {
					sendError(u.WeComBotKey, err, &u)
				}
				// fallthrough
			set:
				needSave = true
				cfg.SetUser(k, u)
			}
		}

		if needSave {
			cfg.Save()
		}

		<-ticker.C
		// }
	}
}

func checkBalanceLoop() {
	ticker := time.NewTicker(time.Duration(cfg.CheckBalanceInterval) * time.Second)
	for {
		// select {
		// case <-ticker.C:
		needSave := false
		for k := range cfg.Users {
			u, _ := cfg.GetUser(k) // ensure locking
			if u.Enabled {
				s, err := xfb.GetCardMoney(u.SessionId, u.YmUserId)
				if err != nil {
					slog.Error("unable to query card balance", "err", err, "name", u.Name)
					goto fail
				}
				if s == "- - -" {
					slog.Info(`GetCardMoney returned "- - -"`)
					continue
				}

				{ // make goto work
					balance, err := strconv.ParseFloat(s, 64)
					if err != nil {
						slog.Error("unable to parse card balance", "err", err, "name", u.Name, "rawbalance", s)
						goto fail
					}
					slog.Info("check balance", "name", u.Name, "balance", balance, "threshold", u.Threshold)
					// fmt.Printf("%s, current: %.2f, threshold: %.2f\n", u.Name, balance, u.Threshold)
					err = rechargeToThreshold(balance, &u)
					if err != nil {
						slog.Error("unable to recharge card balance", "err", err, "name", u.Name, "balance", balance)
						goto fail
					}
				}

				// success?

				if u.Failed != 0 {
					u.Failed = 0
					goto set
				}
				continue
			fail:
				u.Failed++
				if u.Failed <= 3 {
					sendError(u.WeComBotKey, err, &u)
				}
				// fallthrough
			set:
				needSave = true
				cfg.SetUser(k, u)
			}
		}

		if needSave {
			cfg.Save()
		}
		// case <-stop:
		// 	return
		// }

		<-ticker.C
	}
}

func main() {
	dbginfo, _ := debug.ReadBuildInfo()
	println(dbginfo.String())
	slog.Warn("Program started")

	cfg = xfbbroker.LoadConfig()

	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(gorad.NewTextFileSlogHandler(cfg.LogFileName, level)))
	// stop := make(chan bool)

	go checkBalanceLoop()
	go checkTransLoop()

	// http.ListenAndServeTLS(config.ListenAddr, "cert.pem", "key.pem", nil)
	http.ListenAndServe(cfg.ListenAddr, xfbbroker.CreateApiServer(cfg))
	fmt.Println("Hello, World!")
}
