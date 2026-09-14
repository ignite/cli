package ed25519_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
)

func TestVerify(t *testing.T) {
	publicKey, _ := hex.DecodeString("0D853FA898A07EB91F618BB3E8B738B0E45BE9B3053799A2C42F8204F5FA3505")
	signature, _ := hex.DecodeString("2B39638983858715AD2FA059665ADFE267936B8F20C4DA01E9650958E0CA65C0255C75164360F468087FE8385140E48EE3471E332472A50AEE599F9D0EADD106")
	if !ed25519.Verify(publicKey, []byte("hello gno.land"), signature) {
		t.Error("verify failed")
	}
}
