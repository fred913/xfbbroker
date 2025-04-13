package xfb

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

const (
	qrCodeEndpointUrl   = "/card/getQRCode"
	qrResultEndpointUrl = "/card/getQRCodeResult"
)

type QrPayCode struct {
	QRCode    string
	SessionID string
	client    *http.Client
	Creation  int64
}

func GenerateQrPayCode(sessionId string) (*QrPayCode, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	form := url.Values{
		"platform":   []string{"WECHAT_H5"},
		"schoolCode": []string{"20090820"}}

	req, err := http.NewRequest("POST", XfbWebApp+qrCodeEndpointUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "shiroJID", Value: sessionId})
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyCnt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d; response body: %s", resp.StatusCode, string(bodyCnt))
	}

	var result struct {
		Success    bool   `json:"success"`
		StatusCode int    `json:"statusCode"`
		Data       string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success || result.StatusCode != 0 {
		return nil, fmt.Errorf("API error: success=%v, statusCode=%d",
			result.Success, result.StatusCode)
	}

	return &QrPayCode{
		QRCode:    result.Data,
		SessionID: sessionId,
		client:    client,
		Creation:  time.Now().Unix(),
	}, nil
}

func (q *QrPayCode) SaveQrImage(filename string) error {
	qr, err := qrcode.New(q.QRCode, qrcode.Medium)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, qr.Image(256))
}

func (q *QrPayCode) GetQrPngBuf(size int) ([]byte, error) {
	qr, err := qrcode.New(q.QRCode, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	return qr.PNG(size)
}

func (q *QrPayCode) GetResult() (map[string]interface{}, error) {
	form := url.Values{}
	form.Add("qrCode", q.QRCode)

	req, err := http.NewRequest("POST", XfbWebApp+qrResultEndpointUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "shiroJID", Value: q.SessionID})

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: program <session-id>")
// 		os.Exit(1)
// 	}

// 	qr, err := GenerateQrPayCode(os.Args[1])
// 	if err != nil {
// 		fmt.Printf("Error generating QR code: %v\n", err)
// 		os.Exit(1)
// 	}

// 	if err := qr.SaveQrImage("qrcode.png"); err != nil {
// 		fmt.Printf("Error saving QR code: %v\n", err)
// 		os.Exit(1)
// 	}

// 	timeout := time.After(totalCheckTime)
// 	ticker := time.NewTicker(checkInterval)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			result, err := qr.GetResult()
// 			if err != nil {
// 				fmt.Printf("Error checking result: %v\n", err)
// 				continue
// 			}

// 			fmt.Println("Check result:", result)
// 			if success, ok := result["success"].(bool); ok && success {
// 				fmt.Println("Payment successful!")
// 				return
// 			}

// 		case <-timeout:
// 			fmt.Println("Timeout waiting for payment confirmation")
// 			return
// 		}
// 	}
// }
