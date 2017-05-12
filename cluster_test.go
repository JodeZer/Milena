package Milena

import "testing"

func TestRemoveCommonTopics(t *testing.T) {
	ls := []string{"a", "b", "c"}
	rs := []string{"a"}
	t.Logf("%v", removeCommonTopics(ls, rs))
}
