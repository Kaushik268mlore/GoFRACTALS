package main

import (
	"time"
)

var randState = uint64(time.Now().UnixNano())

func RandINT64() uint64 {
	randState = (randState ^ (randState << 13)) ^ (randState >> 7) ^ (randState << 17)
	return randState
}
func RandFLOAT() float64 {
	return float64(RandINT64()/2) / (1 << 63)
}
