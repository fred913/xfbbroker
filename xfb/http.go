package xfb

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/yiffyi/gorad/radhttp"
)

const XfbPay = "https://pay.xiaofubao.com"
const XfbWebApp = "https://webapp.xiaofubao.com"
const XfbApp = "https://application.xiaofubao.com"

var client = &http.Client{
	Timeout:       time.Second * 30,
	CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
}

// func Get(url string, sessionId string, v XfbBaseResponse) (newSessionId string, err error) {
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	if len(sessionId) > 0 {
// 		req.AddCookie(&http.Cookie{Name: "shiroJID", Value: sessionId})
// 	}

// 	resp, b, err := radhttp.JSONDo(client, req, v)
// 	if resp == nil {
// 		return
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		err = fmt.Errorf("bad HTTP Status: %s\n\t%s", resp.Status, string(b))
// 		return
// 	}

// 	newSessionId = ""
// 	for _, v := range resp.Cookies() {
// 		if v.Name == "shiroJID" {
// 			newSessionId = v.Value
// 		}
// 	}

// 	if code := v.GetStatusCode(); code == 0 {
// 		err = nil
// 		return
// 	} else {
// 		err = fmt.Errorf("bad statusCode from xfb: %d", code)
// 		return
// 	}
// }

func Post(url string, sessionId string, payload map[string]interface{}, v XfbBaseResponse) (newSessionId string, err error) {
	req, err := radhttp.NewJSONPostRequest(url, payload)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Referer", "https://webapp.xiaofubao.com/card/card_home.shtml?platform=WECHAT_H5&schoolCode=20090820&thirdAppid=wx8fddf03d92fd6fa9")
	// req.Header.Set("Origin", "https://webapp.xiaofubao.com")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.57(0x18003921) NetType/WIFI Language/en")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	// req.Header.Set("Accept", "application/json, text/plain, */*")
	if len(sessionId) > 0 {
		req.AddCookie(&http.Cookie{Name: "shiroJID", Value: sessionId})
	}

	resp, b, err := radhttp.JSONDo(client, req, v)
	if resp == nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad HTTP Status: %s\n\t%s", resp.Status, string(b))
		return
	}

	newSessionId = ""
	for _, v := range resp.Cookies() {
		if v.Name == "shiroJID" {
			newSessionId = v.Value
		}
	}

	if code := v.GetStatusCode(); code == 0 {
		err = nil
		return
	} else {
		err = fmt.Errorf("bad statusCode from xfb: %d", code)
		return
	}
}

func GetRedirectLocation(url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", errors.New("no Location found")
	}

	return loc, nil
}
