package utils

import "strconv"

func ConvertToBin(n int64) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(int(lsb)) + result
	}
	return result
}

func NextPowOf2(cap int64) int64 {
	if cap < 2 {
		return 2
	}
	if cap&(cap-1) == 0 {
		return cap
	}
	cap |= cap >> 1
	cap |= cap >> 2
	cap |= cap >> 4
	cap |= cap >> 8
	cap |= cap >> 16
	return cap + 1
}