package main

import (
	"fmt"
	"time"

	"github.com/shadas/timewheel"
)

func main() {
	tw := timewheel.NewSingleTimeWheel(time.Millisecond*50, 10)
	tw.Run()
	go func() {
		time.Sleep(3 * time.Second)
		tw.AddTimerTask(3*time.Second, func() {
			fmt.Println("xxx 1")
		}, "xxx1")
		tw.AddTimerTask(3*time.Second, func() {
			fmt.Println("xxx 2")
		}, "xxx2")
		time.Sleep(2 * time.Second)
		tw.RemoveTimerTask("xxx1")
	}()
	select {}
}
