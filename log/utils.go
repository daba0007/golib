package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

type logLevel uint16

const (
	// UNKNOWN UNKNOWN level
	UNKNOWN logLevel = iota
	// DEBUG DEBUG level
	DEBUG
	// INFO INFO level
	INFO
	// WARNING WARNING level
	WARNING
	// ERROR ERROR level
	ERROR
	// FATAL FATAL level
	FATAL
)

type logMsg struct {
	level     string
	msg       string
	funcName  string
	fileName  string
	timestamp string
	lineNo    int
}

// Log log struct
type Log struct {
	level       logLevel     // 日志级别
	filePath    string       // 文件路径
	fileDays    string       // 记录文件日期
	fileObj     *os.File     // 文件对象
	maxFileSize int64        // 最大文件大小, 默认256Mb
	saveDays    int          // 文件保存天数
	consoleFlag bool         // 是否在命令行打印日志，默认打印
	textFlag    bool         // 是否记录到文件，默认记录
	logNum      int64        // 当日文件数，默认为1
	logChan     chan *logMsg // 日志文件缓冲区，支持异步日志写入
}

func parseLogLevel(level string) (logLevel, error) {
	s := strings.ToUpper(level)
	switch s {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return UNKNOWN, fmt.Errorf("unknown level")
	}
}

func unparseLogLevel(lv logLevel) (string, error) {
	switch lv {
	case DEBUG:
		return "DEBUG", nil
	case INFO:
		return "INFO", nil
	case WARNING:
		return "WARNING", nil
	case ERROR:
		return "ERROR", nil
	case FATAL:
		return "FATAL", nil
	default:
		return "UNKNOWN", fmt.Errorf("unknown level")
	}
}

func getInfo(n int) (string, string, int) {
	pc, file, line, ok := runtime.Caller(n)
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return "", "", 0
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file)
	return funcName, fileName, line
}
