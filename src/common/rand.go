package common

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// GenerateRequestID 生成一个随机的RequestID
func GenerateRequestID() string {
	return strconv.FormatUint(rand.Uint64(), 10)
}
