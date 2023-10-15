package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	prefix = strings.Trim(prefix, " [],./:\\")

	//only allow 6 chars
	if len(prefix) > 6 {
		prefix = prefix[:6]
	}

	//add space until 6 chars
	if prefix != "" {
		for len(prefix) < 6 {
			prefix += " "
		}
	}

	prefix = "[" + prefix + "]"
	return &Logger{prefix: prefix}
}

func getCallerMethodName(callerStackOffset int) (string, error) {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(5+callerStackOffset, fpcs)
	if n == 0 {
		return "", fmt.Errorf("no caller found")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return "", fmt.Errorf("caller is nil")
	}

	return caller.Name(), nil
}

func getLogPrefixes(message string, callerStackOffset int) (retMes string, timeStr string, runtimeInfo string, methodInfo string) {
	t := time.Now()

	if _, filename, line, isRuntimeInfoOk := runtime.Caller(3 + callerStackOffset); isRuntimeInfoOk {
		runtimeInfo = filepath.Base(filename) + ":" + strconv.Itoa(line) + ":"

		if methodName, err := getCallerMethodName(callerStackOffset); err == nil {
			methodInfo = "(in " + methodName + "):"
		}
	}

	timeStr = t.Format(time.StampMilli)
	retMes = strings.TrimRight(message, "\n")

	return retMes, timeStr, runtimeInfo, methodInfo
}

func (l *Logger) Info(args ...interface{}) {
	_info(l.prefix, args...)
}

// Debug - Log Debug message
func (l *Logger) Debug(v ...interface{}) {
	_debug(l.prefix, v...)
}

// Warning - Log Warning message
func (l *Logger) Warning(v ...interface{}) {
	_warning(l.prefix, v...)
}

// Trace - Log Trace message
func (l *Logger) Trace(v ...interface{}) {
	_trace(l.prefix, v...)
}

// Error - Log Error message
func (l *Logger) Error(v ...interface{}) {
	_error(l.prefix, 0, v...)
}

func _info(name string, v ...interface{}) {
	mes, timeStr, _, _ := getLogPrefixes(fmt.Sprint(v...), 0)
	fmt.Println(timeStr, name, mes)
}

func _debug(name string, v ...interface{}) {
	mes, timeStr, runtimeInfo, _ := getLogPrefixes(fmt.Sprint(v...), 0)
	fmt.Println(timeStr, name, "DEBUG", runtimeInfo, mes)
}

func _warning(name string, v ...interface{}) {
	mes, timeStr, runtimeInfo, _ := getLogPrefixes(fmt.Sprint(v...), 0)
	fmt.Println(timeStr, name, "WARNING", runtimeInfo, mes)
}

func _trace(name string, v ...interface{}) {
	mes, timeStr, runtimeInfo, methodInfo := getLogPrefixes(fmt.Sprint(v...), 0)
	fmt.Println(timeStr, name, "TRACE", runtimeInfo+methodInfo, mes)
}

func _error(name string, callerStackOffset int, v ...interface{}) {
	mes, timeStr, runtimeInfo, methodInfo := getLogPrefixes(fmt.Sprint(v...), callerStackOffset)
	fmt.Println(timeStr, name, "ERROR", runtimeInfo+methodInfo, mes)
}
