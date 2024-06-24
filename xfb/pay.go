package xfb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

func GetUserById(token, ymId string) (sessionId string, data map[string]interface{}, err error) {
	r, err := PostJSON(XfbWebApp+"/user/getUserById", "", map[string]interface{}{
		"platform": "WECHAT_H5",
		"token":    token,
		"ymId":     ymId,
	})
	if err != nil {
		return "", nil, err
	}
	sessionId = r.SessionId
	data = r.Data.(map[string]interface{})
	return
}

func GetCardMoney(sessionId, ymId string) (string, error) {
	r, err := PostJSON(XfbWebApp+"/card/getCardMoney", sessionId, map[string]interface{}{
		"ymId": ymId,
	})
	if err != nil {
		return "", err
	}

	val := r.Data.(string)
	// what's wrong with you?
	if val == "- - -" {
		slog.Debug(`GetCardMoney: "- - -" received`, "body", r)
	}

	return val, err
}

func CardQuerynoPage(sessionId, ymId string, queryTime time.Time) (total int, rows []Trans, err error) {
	b, _, err := sendPost(XfbWebApp+"/routeauth/auth/route/user/cardQuerynoPage", sessionId, map[string]interface{}{
		"queryTime": queryTime.Format("20060102"),
		"ymId":      ymId,
	})
	if err != nil {
		// TODO err
		return
	}

	var r XfbQueryTransResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal body: %w", err)
		return
	}

	if r.StatusCode == 0 {
		err = nil
		total = r.Total
		rows = r.Rows
	} else {
		err = fmt.Errorf("bad statusCode %d from xfb", r.StatusCode)
	}
	return
}

func RechargeOnCard(money, openId, sessionId, ymId string) (string, error) {
	r, err := PostJSON(XfbWebApp+"/order/rechargeOnCardByParam", sessionId, map[string]interface{}{
		"openid":         openId,
		"totalMoney":     money,
		"orderRealMoney": money,
		"rechargeType":   1,
		"subappid":       "wx8fddf03d92fd6fa9",
		"schoolCode":     "20090820",
		"platform":       "WECHAT_H5",
		"sessionId":      sessionId,
		"ymId":           ymId,
	})
	if err != nil {
		return "", err
	}
	return r.Data.(string), err
}

func SignPayCheck(tranNo string) (string, error) {
	r, err := PostJSON(XfbPay+"/pay/sign/signPayCheck", "", map[string]interface{}{
		"tranNo":  tranNo,
		"payType": "WXPAY",
	})
	return r.Message, err
}

func GetSignUrl(tranNo string) (applyId string, jumpUrl string, err error) {
	r, err := PostJSON(XfbPay+"/h5/pay/sign/getSignUrl", "", map[string]interface{}{
		"payType":     "WXPAY",
		"tranNo":      tranNo,
		"signCashier": 0,
	})
	if err != nil {
		return "", "", err
	}

	d := r.Data.(map[string]interface{})
	applyId = d["applyId"].(string)
	jumpUrl = d["jumpUrl"].(string)
	return
}

func QuerySignApplyById(applyId string) (int, error) {
	r, err := PostJSON(XfbPay+"/h5/pay/sign/querySignApplyById", "", map[string]interface{}{
		"applyId": applyId,
	})
	if err != nil {
		return 0, err
	}

	s := r.Data.(map[string]interface{})["status"].(int)
	if s == 3 {
		return s, nil
	}
	if s == 4 {
		return s, errors.New("sign application failed")
	}
	// s == 1 still applying
	return s, fmt.Errorf("unknown status: %d", s)
}

func PayChoose(tranNo string) error {
	_, err := PostJSON(XfbPay+"/pay/unified/choose.shtml", "", map[string]interface{}{
		"tranNo":    tranNo,
		"payType":   "WXPAY",
		"bussiCode": "WXSIGN",
	})
	return err
}

func DoPay(tranNo string) error {
	_, err := PostJSON(XfbPay+"/pay/doPay", "", map[string]interface{}{
		"tranNo": tranNo,
	})
	return err
}
