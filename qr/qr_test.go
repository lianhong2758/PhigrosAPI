package qr_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lianhong2758/PhigrosAPI/qr"
	"github.com/skip2/go-qrcode"
)

func TestQR(t *testing.T) {
	defer  os.Remove("qr.png")
	r, err := qr.LoginQrCode(true, "public_profile")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("QR_Url:", r.Data.QrcodeURL)
	var png []byte
	png, _ = qrcode.Encode(r.Data.QrcodeURL, qrcode.Medium, 256)
	_ = os.WriteFile("qr.png", png, 0666)
	//重复检查扫码结果
	result, _ := qr.CheckQRCode(true, r)
	//假设五分钟超时
	deadline := time.Now().Add(5 * time.Minute)
	for !result.Success {
		time.Sleep(2 * time.Second)
		result, err = qr.CheckQRCode(true, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		if time.Now().After(deadline) {
			fmt.Println("登录超时!")
			return
		}
	}
	p, err := qr.GetProfile(true, result)
	if err != nil {
		fmt.Println(err)
		return
	}
	k, err := qr.LoginAndGetToken(result, p, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("SessionToken", k.SessionToken)
}
