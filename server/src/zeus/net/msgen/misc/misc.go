package misc

import "sort"

// OrOrderedForEach 对Map元素按键值排序，然后依次调用f(key, value).
func OrderedForEach(input map[string]string, f func(key, value interface{}) bool) {
	// 必须排序，golang中的map都是乱序的
	var sortLst []string
	for k := range input {
		sortLst = append(sortLst, k)
	}
	sort.Strings(sortLst)

	for _, idxVal := range sortLst {
		if val, ok := input[idxVal]; ok {
			if f(idxVal, val) != true {
				break
			}
		}
	}
}
