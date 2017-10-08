package main

import (
	"canihazmidi/src/prolink2midi/converter"
	"canihazmidi/src/prolink2midi/converter/midi"
	"time"
)

func main() {

	conv := converter.NewConverter()
	_ = conv

	time.Sleep(10 * time.Minute)

}
