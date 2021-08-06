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
