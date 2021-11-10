package vietqr

import (
	"fmt"
)

type VietQRError struct {
	Code string
	Desc string
}

func (e VietQRError) Error() string {
	return fmt.Sprintf("VietQR(%s): %s", e.Code, e.Desc)
}

func IsTokenError(err error) bool {
	if e, ok := err.(VietQRError); ok {
		return e.Code == "05"
	}

	return false
}

type Envelope struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Data string `json:"data"`
}

func (e Envelope) Err() error {
	if e.Code != "00" {
		return VietQRError{
			Code: e.Code,
			Desc: e.Desc,
		}
	}

	return nil
}
