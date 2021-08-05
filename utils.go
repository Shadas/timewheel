package timewheel

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
	// offset -> pos
	newPos = make([]int, len(pos))
	addPos := make([]int, len(pos))
	for i := 0; i < len(pos); i++ {
		addPos[i] = offset % multipleLevelSize[MultipleLevel(i+1)]
		offset = offset / multipleLevelSize[MultipleLevel(i+1)]
	}

	more := false
	for i := 0; i < len(pos); i++ {
		sum := addPos[i] + pos[i]
		if more {
			sum += 1
		}
		if sum/multipleLevelSize[MultipleLevel(i+1)] > 0 {
			more = true
		} else {
			more = false
		}
		newPos[i] = sum % multipleLevelSize[MultipleLevel(i+1)]
	}
	return
}
