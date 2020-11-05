package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logger *Log

func init() {
	NewLog("info", "")
}

// NewLog init new console log
func NewLog(l string, filePath string) *Log {
	var (
		err   error
		level logLevel
	)
	if level, err = parseLogLevel(l); err != nil {
		panic(err)
	}
	logger = &Log{
		level:       level,
		filePath:    filePath,
		maxFileSize: int64(1024 * 1024 * 256),
		saveDays:    7,
		consoleFlag: true,
		textFlag:    true,
		logNum:      1,
		logChan:     make(chan *logMsg, 50000),
	}
	if filePath == "" {
		logger.textFlag = false
		return logger
	}
	if err = logger.initFile(filePath); err != nil {
		panic(err)
	}
	return logger
}

// SetSaveDays 设置日志文件数
func (l *Log) SetSaveDays(days int) {
	l.saveDays = days
}

func (l *Log) initFile(filePath string) error {
	fileObj, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file error")
	}
	l.fileObj = fileObj
	// 开启后台goroutine写日志
	go l.writeLogBackground()
	return nil
}

// Close close file
func (l *Log) Close() {
	l.fileObj.Close()
}

// SetLevel set log level
func (l *Log) SetLevel(lv string) {
	var (
		err   error
		level logLevel
	)
	if level, err = parseLogLevel(lv); err != nil {
		fmt.Println(err)
	}
	l.level = level
}

// SetConsoleFlag set console print
func (l *Log) SetConsoleFlag(flag bool) {
	l.consoleFlag = flag
}

// SetTextFlag set text print
func (l Log) SetTextFlag(flag bool) {
	l.textFlag = flag
}

func (l *Log) enable(level logLevel) bool {
	return l.level <= level
}

func (l *Log) console(lv logLevel, str string, msg ...interface{}) {
	if l.enable(lv) {
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		level, err := unparseLogLevel(lv)
		if err != nil {
			fmt.Printf("unparse LogLevel failed\n")
		}
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006/01/02 15:04:05"), colorMsg(lv, level), funcName, fileName, lineNo, colorMsg(lv, fmt.Sprintf(str, msg...)))
	}
}

func (l *Log) writeLogBackground() {
	if l.checkSize(l.fileObj) {
		// 文件过大，分文件记录
		l.fileObj.Close()
		nowStr := time.Now().Format("20060102")
		newLogName := fmt.Sprintf("%s.%s.%d.log", strings.TrimSuffix(l.filePath, ".log"), nowStr, l.logNum)
		os.Rename(l.filePath, newLogName)
		l.logNum++
		fileObj, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("open file error\n")
		}
		l.fileObj = fileObj
	}
	if l.checkDays(l.fileObj) {
		// 新一天，换新日志
		l.fileObj.Close()
		l.logNum = 1
		nowStr := time.Now().Format("20060102")
		newLogName := fmt.Sprintf("%s.%s.%d.log", strings.TrimSuffix(l.filePath, ".log"), nowStr, l.logNum)
		os.Rename(l.filePath, newLogName)
		fileObj, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("open file error\n")
		}
		l.fileObj = fileObj

	}
	for {
		select {
		case logTmp := <-l.logChan:
			fmt.Fprintf(l.fileObj, "[%s] [%v] [%s:%s:%d] %s\n", logTmp.timestamp, logTmp.level, logTmp.funcName, logTmp.fileName, logTmp.lineNo, logTmp.msg)
		default:
			// 取不到日志，休息500ms
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (l *Log) text(lv logLevel, str string, msg ...interface{}) {
	if l.enable(lv) {
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		level, err := unparseLogLevel(lv)
		if err != nil {
			fmt.Printf("unparse LogLevel failed\n")
		}
		// 先把日志发送到通道中
		logTmp := &logMsg{
			level:     level,
			msg:       fmt.Sprintf(str, msg...),
			funcName:  funcName,
			fileName:  fileName,
			timestamp: now.Format("2006/01/02 15:04:05"),
			lineNo:    lineNo,
		}
		select {
		case l.logChan <- logTmp:
		default: // 丢掉日志，不出现阻塞
		}

	}
}

func (l *Log) checkSize(file *os.File) bool {
	var (
		err      error
		fileinfo os.FileInfo
	)
	if fileinfo, err = file.Stat(); err != nil {
		fmt.Printf("get file info failed, err:%v\n ", err)
	}
	// 如果当前文件大小大于等于日志文件的最大值，返回true
	if fileinfo.Size() > l.maxFileSize {
		return true
	}
	return false
}

func (l *Log) checkDays(file *os.File) bool {
	// 如果当前文件日期大于等于日志文件的最大值，返回true
	var (
		err      error
		fileinfo os.FileInfo
	)
	if fileinfo, err = file.Stat(); err != nil {
		fmt.Printf("get file info failed, err:%v\n ", err)
	}

	if fileinfo.ModTime().Day() < time.Now().Day() {
		l.delete()
		return true
	}
	return false
}

func (l *Log) delete() {
	filePath, _ := filepath.Abs(filepath.Dir(l.fileObj.Name()))
	fakeNow := time.Now().AddDate(0, 0, -l.saveDays)
	filepath.Walk(filePath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			// 防止误删
			if !f.IsDir() && f.ModTime().Before(fakeNow) && strings.HasSuffix(f.Name(), l.fileObj.Name()) {
				os.Remove(path)
				return nil
			}

			return nil
		})
}

// Debug Debug level
func Debug(str string, msg ...interface{}) {
	if logger.consoleFlag {
		logger.console(DEBUG, str, msg...)
	}
	if logger.textFlag {
		logger.text(DEBUG, str, msg...)
	}
}

// Info Info level
func Info(str string, msg ...interface{}) {
	if logger.consoleFlag {
		logger.console(INFO, str, msg...)
	}
	if logger.textFlag {
		logger.text(INFO, str, msg...)
	}
}

// Warning Warning level
func Warning(str string, msg ...interface{}) {
	if logger.consoleFlag {
		logger.console(WARNING, str, msg...)
	}
	if logger.textFlag {
		logger.text(WARNING, str, msg...)
	}
}

// Error Error level
func Error(str string, msg ...interface{}) {
	if logger.consoleFlag {
		logger.console(ERROR, str, msg...)
	}
	if logger.textFlag {
		logger.text(ERROR, str, msg...)
	}
}

// Fatal Fatal level
func Fatal(str string, msg ...interface{}) {
	if logger.consoleFlag {
		logger.console(FATAL, str, msg...)
	}
	if logger.textFlag {
		logger.text(FATAL, str, msg...)
	}
}
