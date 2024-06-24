package xfbbroker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type WeComResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func SendWeComMsg(msg map[string]interface{}, key string) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("unable to marshal body: %w", err)
	}

	resp, err := http.Post("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+key, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("unable to perform send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad HTTP Status: %s\n\t%s", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read body: %w", err)
	}

	var r WeComResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return fmt.Errorf("unable to unmarshal body: %w", err)
	}

	if r.ErrCode != 0 {
		return errors.New(r.ErrMsg)
	}
	slog.Info("WeCom send", "body", string(b))
	return nil
}
