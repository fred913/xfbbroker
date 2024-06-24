package xfb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const XfbPay = "https://pay.xiaofubao.com"
const XfbWebApp = "https://webapp.xiaofubao.com"
const XfbApp = "https://application.xiaofubao.com"

var client = http.Client{
	Timeout:       time.Second * 30,
	CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
}

func sendPost(url string, sessionId string, payload map[string]interface{}) ([]byte, []*http.Cookie, error) {
	var body, err = json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to marshal body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if len(sessionId) > 0 {
		req.AddCookie(&http.Cookie{Name: "shiroJID", Value: sessionId})
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to perform request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("bad HTTP Status: %s\n\t%s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read body: %w", err)
	}

	return b, resp.Cookies(), nil
}

func PostJSON(url string, sessionId string, payload map[string]interface{}) (*XfbResponse, error) {
	b, cookies, err := sendPost(url, sessionId, payload)
	if err != nil {
		return nil, err
	}

	var r XfbResponse
	r.SessionId = sessionId
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal body: %w", err)
	}

	for _, v := range cookies {
		if v.Name == "shiroJID" {
			r.SessionId = v.Value
		}
	}

	if r.StatusCode == 0 {
		return &r, nil
	} else {
		return &r, fmt.Errorf("bad statusCode %d from xfb: %s", r.StatusCode, r.Message)
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
