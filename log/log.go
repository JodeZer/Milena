package log

import (
	log "github.com/JodeZer/logrus"
	"runtime"
	"strconv"
	"strings"
)

// log instance
var Log *log.Logger

func init() {
	formatter := &log.TextFormatter{}
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	Log = log.New()
	Log.Formatter = formatter
	Log.Level = log.DebugLevel
}

// Errorf log
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Infof log
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warnf log
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Degbugf log
func Degbugf(format string, args ...interface{}) {
	Log.Debugf(getCodeLine()+format, args...)
}

func getCodeLine() string {

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "file miss"
	}
	lineStr := strconv.Itoa(line)
	ss := strings.Split(file, "Milena/")
	return ss[len(ss)-1] + ":" + lineStr + " "
}
