package main

import (
	"math"
	"syscall/js"
	"unsafe"
)

func osc(freq float64) chan float64 {
	out := make(chan float64)
	go func() {
		phase := 0.0
		for {
			out <- math.Sin(phase)
			phase += 2 * math.Pi * freq / sampleRate
		}
	}()
	return out
}

var sampleRate float64
var out chan float64

func setup(_sampleRate float64) {
	sampleRate = _sampleRate
	out = osc(440)
}

func process(output []float32) int {
	for i := range output {
		output[i] = float32(<-out)
	}
	return len(output)
}

func main() {
	js.Global().Set("goSetup", js.FuncOf(func(t js.Value, a []js.Value) interface{} {
		setup(a[0].Float())
		return nil
	}))
	output := make([]float32, 1024)
	js.Global().Set("goProcess", js.FuncOf(func(t js.Value, a []js.Value) interface{} {
		if length := a[0].Length() / 4; length != len(output) {
			output = make([]float32, length)
		}
		length := process(output)
		js.CopyBytesToJS(a[0], unsafe.Slice((*byte)(unsafe.Pointer(&output[0])), length*4))
		return length
	}))
	<-make(chan int)
}
