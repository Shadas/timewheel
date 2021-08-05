package timewheel

import "container/list"

// timerVec为定时器载体
// index 为当前slot的下标
// slots 为链表数组，挂载
type timerVec struct {
	idx   int
	slots []*list.List
}
