package utils

import ("testing"
	otp "src/utils"
)
func TestOtp(t *testing.T){

	nbytes := 8
	iterations := 100
	for i:=0;i<iterations;i++{
		code , _ := otp.GenerateRandomOTP(nbytes)
		t.Log(code)
	}
}