package utils

import (
	"math/big"
	"time"
)

type TokenIDReader struct {
	f func() int64
}

var tr *TokenIDReader

func init() {
	tr = NewTokenIDReader(NextID())
}

func NewTokenIDReader(f func() int64) *TokenIDReader {
	return  &TokenIDReader{f}
}

// Closure
func NextID() func() int64 {
	var index int64 = 0
	return func() int64 {
		if index >= 10000 {
			index = 0
		}
		index++
		return index
	}
}

func NewTokenID() int64 {
	value1 := big.NewInt(tr.f())
	value2 := big.NewInt(0)
	value2, _ = value2.SetString(time.Now().Format("20060102150405"), 10)
	return value2.Int64() * 10000 + value1.Int64()
}