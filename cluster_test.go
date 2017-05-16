package Milena

import "testing"

func TestRemoveCommonTopics(t *testing.T) {
	ls := []string{"a","a", "b","b", "c"}
	rs := []string{"a"}
	res := removeCommonTopics(ls, rs)
	t.Logf("%v %d", res,len(res))
}
