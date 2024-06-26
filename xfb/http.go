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

func Post(url string, sessionId string, payload map[string]interface{}, v XfbBaseResponse) (newSessionId string, err error) {
	req, err := radhttp.NewJSONPostRequest(url, payload)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
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
