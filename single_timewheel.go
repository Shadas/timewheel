package timewheel

import (
	"container/list"
	"fmt"
	"time"
)

type Task struct {
	delay time.Duration
	f     func()
	key   interface{} // 用于删除定时任务

	circle int
}

// 单层时间轮
type SingleTimeWheel struct {
	baseTimeWheel
	totalDuration time.Duration
	timerVec      // 对应定时器载体

	keySlot map[interface{}]int // 记录key-slot位置，用于删除
}

func NewSingleTimeWheel(interval time.Duration, slotNum int) *SingleTimeWheel {
	tw := &SingleTimeWheel{
		baseTimeWheel: baseTimeWheel{
			interval:    interval,
			slotNum:     slotNum,
			shutdown:    make(chan struct{}),
			addTaskChan: make(chan *Task),
			rmTaskChan:  make(chan interface{}),
		},
		timerVec: timerVec{
			slots: make([]*list.List, slotNum),
		},
		totalDuration: time.Duration(int64(interval) * int64(slotNum)),
		keySlot:       make(map[interface{}]int),
	}
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
	return tw
}

// 时间轮运行
func (tw *SingleTimeWheel) Run() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

func (tw *SingleTimeWheel) run() {
	for {
		select {
		case <-tw.ticker.C: // 指针滴答
			tw.processTick()
		case task := <-tw.addTaskChan: // 增加新任务
			tw.addTask(task)
		case key := <-tw.rmTaskChan: // 删除任务
			tw.removeTask(key)
		case <-tw.shutdown: // 停止
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *SingleTimeWheel) RemoveTimerTask(key interface{}) {
	if key == nil {
		return
	}
	tw.rmTaskChan <- key
}

func (tw *SingleTimeWheel) removeTask(key interface{}) {
	idx, ok := tw.keySlot[key]
	if !ok {
		return
	}
	l := tw.slots[idx]
	tw.removeTaskFromList(l, key)
}

func (tw *SingleTimeWheel) removeTaskFromList(l *list.List, key interface{}) {
	for node := l.Front(); node != nil; {
		task := node.Value.(*Task)
		if task.key == key {
			next := node.Next()
			l.Remove(node)
			node = next
		} else {
			node = node.Next()
		}
	}
}

func (tw *SingleTimeWheel) AddTimerTask(delay time.Duration, f func(), key interface{}) {
	if delay < 0 {
		return
	}
	tw.addTaskChan <- &Task{
		delay: delay,
		f:     f,
		key:   key,
	}
}

func (tw *SingleTimeWheel) addTask(task *Task) {
	fmt.Printf("add Task, task=%+v\n", *task)
	idx := tw.parseIdxAndCircle(task)
	tw.slots[idx].PushBack(task)
	if task.key != nil {
		tw.keySlot[task.key] = idx
	}
}

func (tw *SingleTimeWheel) parseIdxAndCircle(task *Task) (idx int) {
	task.circle = int(task.delay / tw.totalDuration)
	offset := int(task.delay/tw.interval) % tw.slotNum
	idx = (tw.idx + offset) % tw.slotNum
	fmt.Printf("circle=%d, offset=%d, idx=%d\n", task.circle, offset, idx)
	return
}

func (tw *SingleTimeWheel) processTick() {
	fmt.Printf("currPos=%d\n", tw.idx)
	l := tw.slots[tw.idx]
	tw.processList(l)
	if tw.idx == tw.slotNum-1 {
		tw.idx = 0
	} else {
		tw.idx += 1
	}
}

func (tw *SingleTimeWheel) processList(l *list.List) {
	for node := l.Front(); node != nil; {
		task := node.Value.(*Task)
		fmt.Println("curr_circle=", task.circle)
		if task.circle > 0 {
			task.circle -= 1
			node = node.Next()
			continue
		}
		go task.f()
		next := node.Next()
		l.Remove(node)
		node = next
	}
}
