package array

import "strings"

func ContainsStr(array []string, expected string) bool {
	for _, value := range array {
		if strings.Contains(expected, value) {
			return true
		}
	}
	return false
}
