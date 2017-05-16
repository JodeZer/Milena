package Milena

import (
	"testing"
	"sort"
	"github.com/JodeZer/Milena/util/testUtil"
)

type sortStringSlice []string

func (s sortStringSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortStringSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortStringSlice) Len() int {
	return len(s)
}

func TestRemoveCommonTopics(t *testing.T) {

	f := func(a interface{}, b interface{}) bool {
		as := a.(sortStringSlice)
		bs := b.(sortStringSlice)
		if len(as) != len(bs) {
			return false
		}
		sort.Sort(as)
		sort.Sort(bs)
		for i, v := range as {
			if v != bs[i] {
				return false
			}
		}
		return true
	}
	ls := sortStringSlice([]string{"a", "a", "b", "b", "c"})
	rs := sortStringSlice([]string{"a"})
	testUtil.Test(t, "TestRemoveCommonTopics1", sortStringSlice(removeCommonTopics(ls, rs)), sortStringSlice([]string{"b", "c"}), f)

	ls = sortStringSlice([]string{"a", "a", "b", "b", "c"})
	rs = sortStringSlice([]string{"d"})
	testUtil.Test(t, "TestRemoveCommonTopics2", sortStringSlice(removeCommonTopics(ls, rs)), sortStringSlice([]string{"b", "c", "a"}), f)

	ls = sortStringSlice([]string{})
	rs = sortStringSlice([]string{"d"})
	testUtil.Test(t, "TestRemoveCommonTopics3", sortStringSlice(removeCommonTopics(ls, rs)), sortStringSlice([]string{}), f)

}
