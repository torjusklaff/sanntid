package main

import (
	. "fmt"
	"runtime"
	"time"
)

var j int = 0
var mutex = make(chan int, 1)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
	// Try doing the exercise both with and without it!

	go someGoroutine1()
	go someGoroutine2() // This spawns someGoroutine() as a goroutine

	// We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
	// We'll come back to using channels in Exercise 2. For now: Sleep.
	time.Sleep(100 * time.Millisecond)
	Println(j)
}
func someGoroutine2() {
	for i := 0; i < 1000000; i++ {
		mutex <- 1
		j++
		<-mutex
	}
}
func someGoroutine1() {
	for i := 0; i < 1000000; i++ {
		mutex <- 1
		j--
		<-mutex
	}

}

