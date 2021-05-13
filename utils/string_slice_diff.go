package utils

import (
	"sort"
	"strings"
)

func StringSliceDiff(a, b []string) []string {
	a = sortIfNeeded(a)
	b = sortIfNeeded(b)
	var d []string
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		c := strings.Compare(a[i], b[j])
		if c == 0 {
			i++
			j++
		} else if c < 0 {
			d = append(d, a[i])
			i++
		} else {
			d = append(d, b[j])
			j++
		}
	}
	d = append(d, a[i:]...)
	d = append(d, b[j:]...)
	return d
}

func sortIfNeeded(a []string) []string {
	if sort.StringsAreSorted(a) {
		return a
	}
	s := append(a[:0:0], a...)
	sort.Strings(s)
	return s
}
