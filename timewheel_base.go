package timewheel

import "time"

type BaseTimeWheel struct {
	interval time.Duration // 时间轮指针移动间隔
	ticker   *time.Ticker
	slotNum  int64

	shutdown    chan struct{}
	addTaskChan chan *Task
	rmTaskChan  chan interface{}
}
