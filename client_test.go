package vietqr_test

import (
	"context"
	"testing"

	vietqr "giautm.dev/vietqr"
)

func TestClient_GetBankList(t *testing.T) {
	type fields struct {
		http vietqr.HttpClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: context.Background(),
			},
			want: 52,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vietqr.NewClient(tt.fields.http)
			c.RequestAccessToken(tt.args.ctx)

			got, err := c.GetBankList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetBankList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Client.GetBankList() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestClient_GenQRCode(t *testing.T) {
	type fields struct {
		http vietqr.HttpClient
	}
	type args struct {
		ctx   context.Context
		input vietqr.GenQRCodeInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: context.Background(),
				input: vietqr.GenQRCodeInput{
					IsMask:      1,
					AcqID:       "970423",
					AccountNo:   "06202945202", // 19
					AccountName: "A B C",       // 50
					Amount:      "100000",      // 13
					Message:     "noi dung",    // 25
				},
			},
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vietqr.NewClient(tt.fields.http)
			c.RequestAccessToken(tt.args.ctx)

			got, err := c.GenQRCode(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GenQRCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil || len(got.ImagePNG) == 0 {
				t.Errorf("Client.GenQRCode() = %v", got)
			}
		})
	}
}
