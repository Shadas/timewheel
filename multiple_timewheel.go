package timewheel

import (
	"container/list"
	"fmt"
	"time"
)

type MultipleLevel int

const (
	MultipleLevelNull MultipleLevel = iota
	MultipleLevel1
	MultipleLevel2
	MultipleLevel3
	MultipleLevel4
	MultipleLevel5
)

const (
	tvrBits = 8
	tvnBits = 6
	tvrSize = 1 << tvrBits
	tvnSize = 1 << tvnBits
)

var multipleLevelSize = map[MultipleLevel]int{
	MultipleLevel1: tvrSize,
	MultipleLevel2: tvnSize,
	MultipleLevel3: tvnSize,
	MultipleLevel4: tvnSize,
	MultipleLevel5: tvnSize,
}

var multipleLevelBase = map[MultipleLevel]int{
	MultipleLevel1: 1,
	MultipleLevel2: 1 << tvrBits,
	MultipleLevel3: 1 << tvrBits << tvnBits,
	MultipleLevel4: 1 << tvrBits << tvnBits << tvnBits,
	MultipleLevel5: 1 << tvrBits << tvnBits << tvnBits << tvnBits,
}

// 多层时间轮，按照输入slotNum决定层数
type MultipleTimeWheel struct {
	baseTimeWheel
	totalDuration time.Duration
	wheels        []*timerVec   // 按序为每一层的时间轮
	level         MultipleLevel // 时间轮层数
}

// interval 为每个最小单位槽位指针移动的间隔
func NewMultipleTimeWheel(interval time.Duration, level MultipleLevel) *MultipleTimeWheel {
	lv, wheels, slotNum := initMultipleTimeWheels(level)
	tw := &MultipleTimeWheel{
		baseTimeWheel: baseTimeWheel{
			interval:     interval,
			slotNum:      slotNum,
			shutdown:     make(chan struct{}),
			addMTaskChan: make(chan *MTask),
			rmTaskChan:   make(chan interface{}),
		},
		totalDuration: time.Duration(int64(interval) * int64(slotNum)),
		wheels:        wheels,
		level:         lv,
	}
	return tw
}

// 根据时间轮层数决定数据结构分配
func initMultipleTimeWheels(l MultipleLevel) (level MultipleLevel, wheels []*timerVec, slotNum int) {
	if l <= MultipleLevelNull {
		level = MultipleLevel1
	} else if l > MultipleLevel5 {
		level = MultipleLevel5
	} else {
		level = l
	}
	wheels = make([]*timerVec, level)
	for i := 0; i < int(level); i++ {
		tv := &timerVec{}
		if i == 0 {
			tv.slots = make([]*list.List, tvrSize)
			for idx := 0; idx < tvrSize; idx++ {
				tv.slots[idx] = list.New()
			}
			slotNum = tvrSize
		} else {
			tv.slots = make([]*list.List, tvnSize)
			for idx := 0; idx < tvnSize; idx++ {
				tv.slots[idx] = list.New()
			}
			slotNum *= tvnSize
		}
		wheels[i] = tv
	}
	return
}

// 时间轮运行
func (tw *MultipleTimeWheel) Run() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

func (tw *MultipleTimeWheel) run() {
	for {
		select {
		//case <-tw.ticker.C: // 指针滴答
		//	tw.processTick()
		case task := <-tw.addMTaskChan: // 增加新任务
			tw.addTask(task)
		//case key := <-tw.rmTaskChan: // 删除任务
		//	tw.removeTask(key)
		case <-tw.shutdown: // 停止
			tw.ticker.Stop()
			return
		}
	}
}

// 设置时间轮指针位置
func (tw *MultipleTimeWheel) SetCurPos(pos []int) {
	for i := 0; i < int(tw.level); i++ {
		if len(pos) < i+1 {
			tw.wheels[i].idx = 0
		} else {
			tw.wheels[i].idx = pos[i]
		}
	}
}

type MTask struct {
	pos   []int
	delay time.Duration
	f     func()
	key   interface{}

	circle  int
	initPos []int // 用于支持判断circle何时应该-1
}

// 获取当前时间指针位置
func (tw *MultipleTimeWheel) CurPos() (pos []int) {
	pos = make([]int, tw.level)
	for idx, wheel := range tw.wheels {
		pos[idx] = wheel.idx
	}
	return
}

func (tw *MultipleTimeWheel) appendPos(offset int) (pos []int, circle int) {
	circle = offset / tw.slotNum // 总圈数=移动格数/一圈总格数
	last := offset % tw.slotNum
	pos = posAddOffset(tw.CurPos(), last)
	return
}

// 增加定时任务
func (tw *MultipleTimeWheel) AddTimerTask(delay time.Duration, f func(), key interface{}) {
	if delay < 0 {
		return
	}
	tw.addMTaskChan <- &MTask{
		delay: delay,
		f:     f,
		key:   key,
		pos:   make([]int, tw.level),
	}
}

func (tw *MultipleTimeWheel) addTask(task *MTask) {
	fmt.Printf("add Task, task=%+v\n", *task)
	offset := int(task.delay / tw.interval)
	task.pos, task.circle = tw.appendPos(offset)
	task.initPos = posAddOffset(tw.CurPos(), -1) // 注册位置标记为上一个ticker的pos
	tw.placeTask(task)                           // 挂载链表
}

// 挂载任务，最开始的任务挂载到最外层时间轮，后续根据ticker触发向内层转移
func (tw *MultipleTimeWheel) placeTask(task *MTask) {
	high := task.pos[tw.level-1]
	l := tw.wheels[tw.level-1].slots[high]
	l.PushBack(task)
}

// 定时触发逻辑(not finish)
func (tw *MultipleTimeWheel) processTick() {
	fmt.Printf("currPos=%d\n", tw.CurPos())
	curPos := tw.CurPos()
	// 先执行最内圈任务
	l := tw.wheels[0].slots[curPos[0]]
	tw.processList(l)
	// 高层时间轮向内转移
	if tw.level > 1 {
		for i := 1; i < int(tw.level); i++ {
			if i == int(tw.level)-1 { // 如果是最高轮，需要关注circle
				//l := tw.pickMovingTasks() // 获取d
			} else {

			}
		}
	}

}

func (tw *MultipleTimeWheel) processList(l *list.List) {
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
