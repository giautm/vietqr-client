package vietqr

import (
	"context"
	"net/http"
	"time"
)

const BaseURL = "https://vietqr.net"

var (
	TimeNow       = time.Now
	GenPassphrase = func(accessToken string) string {
		return TimeNow().Format("20060102") + accessToken
	}
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	AccessToken string
	BaseURL     string
	http        HttpClient
}

func NewClient(httpClient HttpClient) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		AccessToken: "",
		BaseURL:     BaseURL,
		http:        httpClient,
	}
}

func (c *Client) RequestAccessToken(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodGet, "/portal-service/api/data/generate", nil)
	if err != nil {
		return err
	}

	token, err := c.request(req)
	if err == nil {
		c.AccessToken = token
	}
	return err
}

type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
	BIN  string `json:"bin"`
}

func (c *Client) GetBankList(ctx context.Context) ([]Bank, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/portal-service/v1/info/banks", nil)
	if err != nil {
		return nil, err
	}

	var banks []Bank
	if err = c.requestWithDecrypt(req, &banks); err != nil {
		return nil, err
	}
	return banks, nil
}

type GenQRCodeInput struct {
	IsMask      int    `json:"isMask"`
	AcqID       string `json:"acqId"`
	AccountNo   string `json:"accountNo"`
	AccountName string `json:"accountName"`
	Amount      string `json:"amount"`
	Message     string `json:"addInfo"`
}

type GenQRCodeResult struct {
	ImagePNG []byte `json:"qrBase64"`
}

func (c *Client) GenQRCode(ctx context.Context, input GenQRCodeInput) (*GenQRCodeResult, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/portal-service/v1/qr-ibft/generate", input)
	if err != nil {
		return nil, err
	}

	var result GenQRCodeResult
	if err = c.requestWithDecrypt(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
