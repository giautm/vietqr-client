package vietqr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	openssl "github.com/Luzifer/go-openssl/v4"
)

var o = openssl.New()

func (c *Client) encryptPayload(input interface{}) (io.Reader, string, error) {
	if input == nil {
		return nil, "application/json", nil
	}

	plainData, err := json.Marshal(input)
	if err != nil {
		return nil, "", err
	}
	encryptedBase64Data, err := o.EncryptBytes(GenPassphrase(c.AccessToken),
		plainData, openssl.BytesToKeyMD5)
	if err != nil {
		return nil, "", err
	}

	body := bytes.NewBuffer(nil)
	w := multipart.NewWriter(body)
	w.WriteField("payload", string(encryptedBase64Data))

	err = w.Close()
	if err != nil {
		return nil, "", err
	}

	return body, w.FormDataContentType(), nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	b, contentType, err := c.encryptPayload(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, b)
	if err == nil {
		req.Header.Add("Content-Type", contentType)
	}

	return req, err
}

func (c *Client) request(req *http.Request) (string, error) {
	if c.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	}
	// NOTE(giautm): We need to send origin header to avoid blocking from VietQR.net
	//
	// ```
	// VietQR(99): Internal Error
	// ```
	req.Header.Add("Origin", "https://vietqr.net")

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var env Envelope
	err = json.NewDecoder(res.Body).Decode(&env)
	if err != nil {
		return "", err
	}
	if err = env.Err(); err != nil {
		return "", err
	}

	return env.Data, nil
}

func (c *Client) requestWithDecrypt(req *http.Request, v interface{}) error {
	encryptedBase64Data, err := c.request(req)
	if err != nil {
		return err
	}

	plainData, err := o.DecryptBytes(GenPassphrase(c.AccessToken),
		[]byte(encryptedBase64Data), openssl.BytesToKeyMD5)
	if err != nil {
		return err
	}

	return json.Unmarshal(plainData, v)
}
