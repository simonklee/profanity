package math

// pkg math implements some common util functions

// IntMin returns the smaller of x or y.
func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// IntMax returns the larger of x or y.
func IntMax(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// UintMin returns the smaller of x or y.
func UintMin(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

// UintMax returns the larger of x or y.
func UintMax(x, y uint) uint {
	if x < y {
		return y
	}
	return x
}
