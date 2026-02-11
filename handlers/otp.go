package handlers

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000 // 6 digit
	return strconv.Itoa(otp)
}
