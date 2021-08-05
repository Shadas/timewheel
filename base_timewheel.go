package timewheel

import "time"

type baseTimeWheel struct {
	interval time.Duration // 时间轮指针移动间隔
	ticker   *time.Ticker
	slotNum  int

	shutdown     chan struct{}
	addTaskChan  chan *Task
	addMTaskChan chan *MTask
	rmTaskChan   chan interface{}
}

// 时间轮终止
func (tw *baseTimeWheel) ShutDown() {
	close(tw.shutdown)
	close(tw.addTaskChan)
	close(tw.rmTaskChan)
}
