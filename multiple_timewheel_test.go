package timewheel

import (
	"container/list"
	"testing"
	"time"
)

func TestCurrPos(t *testing.T) {
	tw := NewMultipleTimeWheel(time.Second, MultipleLevel3)
	curPos := tw.CurPos()
	if !isPosEqual(curPos, []int{0, 0, 0}) {
		t.Errorf("curPos is not [0,0,0], it's %v", curPos)
	}
	tw.SetCurPos([]int{9, 6, 3})
	curPos = tw.CurPos()
	if !isPosEqual(curPos, []int{9, 6, 3}) {
		t.Errorf("curPos is not [9,6,3], it's %v", curPos)
	}
	tw.SetCurPos([]int{9, 6, 3, 4})
	curPos = tw.CurPos()
	if !isPosEqual(curPos, []int{9, 6, 3}) {
		t.Errorf("curPos is not [9,6,3], it's %v", curPos)
	}
	tw.SetCurPos([]int{5, 5})
	curPos = tw.CurPos()
	if !isPosEqual(curPos, []int{5, 5, 0}) {
		t.Errorf("curPos is not [5,5,0], it's %v", curPos)
	}
}

func TestPosAddOffset(t *testing.T) {
	if pos := posAddOffset([]int{200, 3, 5}, 100); !isPosEqual(pos, []int{44, 4, 5}) {
		t.Errorf("pos is not [44,4,5], it's %v", pos)
	}
	if pos := posAddOffset([]int{200, 3, 5}, 30); !isPosEqual(pos, []int{230, 3, 5}) {
		t.Errorf("pos is not [230,3,5], it's %v", pos)
	}
	if pos := posAddOffset([]int{200, 63, 5}, 1024); !isPosEqual(pos, []int{200, 3, 6}) {
		t.Errorf("pos is not [200,3,6], it's %v", pos)
	}
	if pos := posAddOffset([]int{200, 63, 5}, -1); !isPosEqual(pos, []int{199, 63, 5}) {
		t.Errorf("pos is not [199,63,5], it's %v", pos)
	}
	if pos := posAddOffset([]int{0, 63, 5}, -1); !isPosEqual(pos, []int{255, 62, 5}) {
		t.Errorf("pos is not [255,62,5], it's %v", pos)
	}
}

func TestAppendPos(t *testing.T) {
	tw := NewMultipleTimeWheel(time.Second, MultipleLevel3)
	tw.SetCurPos([]int{9, 6, 3})
	pos, circle := tw.appendPos(25)
	if !isPosEqual(pos, []int{34, 6, 3}) {
		t.Errorf("curPos is not [34,6,3], it's %v", pos)
	}
	if circle != 0 {
		t.Error("circle should be 0")
	}
}

func TestPickMovingTasks(t *testing.T) {
	oldList := list.New()
	oldList.PushBack(&MTask{
		circle:  2,
		initPos: []int{6, 1, 2},
		key:     "1",
	})
	oldList.PushBack(&MTask{
		circle:  0,
		initPos: []int{6, 1, 2},
		key:     "2",
	})
	b, newList := pickMovingTasks(oldList, []int{6, 1, 2})
	if !b {
		t.Error("b should be true")
	}
	for node := oldList.Front(); node != nil; node = node.Next() {
		task := node.Value.(*MTask)
		t.Log(*task)
	}
	for node := newList.Front(); node != nil; node = node.Next() {
		task := node.Value.(*MTask)
		t.Log(*task)
	}
}
