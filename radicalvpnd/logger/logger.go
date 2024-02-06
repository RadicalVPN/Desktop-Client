package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hectane/go-acl"
)

var (
	globalFileHandle *os.File
	filePath         string
)

type Logger struct {
	prefix string
}

func Init(logFile string) {
	filePath = logFile
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

func log(data ...interface{}) error {

	//init log file
	if globalFileHandle == nil {
		var err error
		globalFileHandle, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create log file: %w", err)
		}

		if runtime.GOOS == "windows" {
			if err := acl.Chmod(filePath, 0644); err != nil {
				os.Remove(filePath) // #nosec G104
				return fmt.Errorf("failed to change windows file permissions: %w", err)
			}
		}
	}

	//normal log to stdout
	fmt.Println(data...)

	//write to file
	if globalFileHandle != nil {
		globalFileHandle.WriteString(fmt.Sprintln(data...)) // #nosec G104
	}

	return nil
}

func _info(name string, v ...interface{}) {
	mes, timeStr, _, _ := getLogPrefixes(fmt.Sprint(v...), 0)
	log(timeStr, name, "INFO", mes) // #nosec G104
}

func _debug(name string, v ...interface{}) {
	mes, timeStr, runtimeInfo, _ := getLogPrefixes(fmt.Sprint(v...), 0)
	log(timeStr, name, "DEBU", runtimeInfo, mes) // #nosec G104
}

func _warning(name string, v ...interface{}) {
	mes, timeStr, runtimeInfo, _ := getLogPrefixes(fmt.Sprint(v...), 0)
	log(timeStr, name, "WARN", runtimeInfo, mes) // #nosec G104
}

func _trace(name string, v ...interface{}) {
	mes, timeStr, runtimeInfo, methodInfo := getLogPrefixes(fmt.Sprint(v...), 0)
	log(timeStr, name, "TRAC", runtimeInfo+methodInfo, mes) // #nosec G104
}

func _error(name string, callerStackOffset int, v ...interface{}) {
	mes, timeStr, runtimeInfo, methodInfo := getLogPrefixes(fmt.Sprint(v...), callerStackOffset)
	log(timeStr, name, "ERRO", runtimeInfo+methodInfo, mes) // #nosec G104
}
