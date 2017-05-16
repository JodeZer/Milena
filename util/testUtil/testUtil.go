package testUtil

import (
	"testing"
	"strings"
	"strconv"
	"runtime"
)

type DeepEqual func(interface{}, interface{}) bool

func Test(t *testing.T, n string, in interface{}, expect interface{}, f DeepEqual) {
	if f(in, expect) {
		t.Logf("%s pass\n", n)
		return
	}
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		t.Fatalf("%s runtime fail\n", n)
	}
	lineStr := strconv.Itoa(line)
	ss := strings.Split(file, "Milena/")
	LINE := ss[len(ss) - 1] + ":" + lineStr + " "
	t.Fatalf("%s %s fail\n", LINE, n)
}
