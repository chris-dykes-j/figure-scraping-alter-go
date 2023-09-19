package main

import (
	"math/rand"
	"time"
)

func sleepShort() {
	randomNumber := rand.Float64()*(4-2) + 2
	time.Sleep(time.Duration(randomNumber) * time.Second)
}

func sleepLong() {
	randomNumber := rand.Float64()*(10-5) + 5
	time.Sleep(time.Duration(randomNumber) * time.Second)
}
