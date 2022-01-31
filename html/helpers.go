package html

import "strconv"

func containsIdHelper(search int, list []string) bool {
	for _, value := range list {
		if value == strconv.Itoa(search) {
			return true
		}
	}
	return false
}

func getIndexOfArrayHelper(search int, list []string) int {
	for i, value := range list {
		if value == strconv.Itoa(search) {
			return i
		}
	}
	return -1
}