package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func threadA(name string, mm chan int, done chan int) {
	for {
		fmt.Printf("Thread: %s\n", name)
		select {
		case <-done:
			fmt.Printf("Quitting %s", name)
			return
		default:
		}
		value := <-mm
		value += 2
		mm <- value
	}
}

func threadB(name string, mm chan int, done chan int) {
	for {
		fmt.Printf("Thread: %s\n", name)
		select {
		case <-done:
			fmt.Printf("Quitting %s", name)
			return
		default:
		}
		value := <-mm
		value -= 2
		mm <- value
	}
}

func main() {
	fmt.Println("Hallo Welt")

	var count uint64
	count = 0
	memorymap := make(chan int)
	done := make(chan int)
	go threadA("threadA", memorymap, done)
	memorymap <- 4711
	go threadB("threadB", memorymap, done)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		done <- 1
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	for {
		count++
		if count%10 == 0 {
			runtime.Gosched()
		}
	}

}
