package timewheel

import "container/list"

func isPosEqual(s1, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func posAddOffset(pos []int, offset int) (newPos []int) {
	// pos => idx num
	ss := 0
	for i := 0; i < len(pos); i++ {
		ss += pos[i] * multipleLevelBase[MultipleLevel(i+1)]
	}
	news := ss + offset
	// idx num => new pos
	newPos = make([]int, len(pos))
	for i := 0; i < len(pos); i++ {
		newPos[i] = news % multipleLevelSize[MultipleLevel(i+1)]
		news = news / multipleLevelSize[MultipleLevel(i+1)]
	}
	return
}

// 选出需要移动的任务链表, 任务类型断言为MTask
// param:
// 	sl 原始任务链表
// 	curPos 当前位置
// return:
//	hasTask 是否有pick出的任务
//  l 选出的任务链表
func pickMovingTasks(sl *list.List, curPos []int) (hasTask bool, l *list.List) {
	l = list.New()
	for node := sl.Front(); node != nil; {
		task := node.Value.(*MTask)
		if !isPosEqual(task.initPos, curPos) { // 如果任务位置都不对，则检查下一个任务
			node = node.Next()
			continue
		}
		if task.circle != 0 {
			node = node.Next()
			task.circle--
			continue
		}
		l.PushBack(task)
		hasTask = true
		next := node.Next()
		sl.Remove(node)
		node = next
	}
	return
}
