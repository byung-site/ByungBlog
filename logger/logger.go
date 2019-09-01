/**
 * byung.cn用的日志包
 * 功能：
 *     1.自定义日志格式
 *     2.日志按天存储
 *     3.3种日志级别
 */

package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FormatCallback func(logger *Logger, level int, message string) string

type Logger struct {
	mu             sync.Mutex
	out            io.Writer
	calldepth      int
	logDir         string
	keepDays       int
	logs           []string
	curLogfile     *os.File
	formatCallback FormatCallback
}

const (
	Lshortfile = 1 << iota
	Llongfile
)

const (
	ErrorLevel = 1 << iota
	InfoLevel
	DebugLevel
)

var std = New(os.Stdout, 4)

func init() {
	std.formatCallback = stdFormat
}

func pathOrFileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func timeout(logger *Logger) {
	exist, _ := pathOrFileExists(logger.logDir)
	if !exist {
		os.MkdirAll(logger.logDir, os.ModePerm)
	}

	t := time.Now()
	date := t.Format("2006-01-02")
	logFilename := logger.logDir + "/" + date + ".log"

	file, err := os.OpenFile(logFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}

	if logger.curLogfile != nil {
		logger.curLogfile.Close()
	}
	logger.curLogfile = file
	logger.SetOutput(file)

	logger.logs = append(logger.logs, date)
	length := len(logger.logs)
	if length == logger.keepDays+1 {
		os.Remove(logger.logDir + "/" + logger.logs[0])
		logger.logs = logger.logs[1:]
	}
}

func New(out io.Writer, calldepth int) (logger *Logger) {
	logger = new(Logger)

	logger.out = out

	if calldepth == 0 {
		calldepth = 3
	}
	logger.calldepth = calldepth
	return logger
}

func (this *Logger) FileAndLine(flags int) (file string, line int, ok bool) {
	_, file, line, ok = runtime.Caller(this.calldepth)
	if (flags & Lshortfile) != 0 {
		items := strings.Split(file, "/")
		short := items[len(items)-1]
		file = short
	}
	return file, line, ok
}

func (this *Logger) Time(format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return time.Now().Format(format)
}

func (this *Logger) SetOutput(w io.Writer) {
	defer this.mu.Unlock()
	this.mu.Lock()
	this.out = w
}

func (this *Logger) Output(log string) error {
	defer this.mu.Unlock()
	this.mu.Lock()
	_, err := this.out.Write([]byte(log))
	return err
}

func (this *Logger) SetFormatCallback(formatCallback FormatCallback) {
	this.formatCallback = formatCallback
}

func (this *Logger) Error(v ...interface{}) {
	var buf bytes.Buffer
	var log string

	message := fmt.Sprintln(v...)
	if this.formatCallback != nil {
		log = this.formatCallback(this, ErrorLevel, message)
	} else {
		buf.WriteString("[error] ")
		buf.WriteString(message)
		log = buf.String()
	}

	this.Output(log)
}

func (this *Logger) Errorf(format string, v ...interface{}) {
	var buf bytes.Buffer
	var log string

	message := fmt.Sprintf(format, v...)
	if this.formatCallback != nil {
		log = this.formatCallback(this, ErrorLevel, message)
	} else {
		buf.WriteString("[error] ")
		buf.WriteString(message)
		log = buf.String()
	}

	this.Output(log)
}

func (this *Logger) Info(v ...interface{}) {
	var buf bytes.Buffer
	var log string

	message := fmt.Sprintln(v...)
	if this.formatCallback != nil {
		log = this.formatCallback(this, InfoLevel, message)
	} else {
		buf.WriteString("[info] ")
		buf.WriteString(message)
		log = buf.String()
	}

	this.Output(log)
}

func (this *Logger) Infof(format string, v ...interface{}) {
	var buf bytes.Buffer
	var log string

	message := fmt.Sprintf(format, v...)
	if this.formatCallback != nil {
		log = this.formatCallback(this, InfoLevel, message)
	} else {
		buf.WriteString("[info] ")
		buf.WriteString(message)
		log = buf.String()
	}

	this.Output(log)
}

func (this *Logger) Debug(v ...interface{}) {
	var buf bytes.Buffer
	var log string

	message := fmt.Sprintln(v...)
	if this.formatCallback != nil {
		log = this.formatCallback(this, DebugLevel, message)
	} else {
		buf.WriteString("[debug] ")
		buf.WriteString(message)
		log = buf.String()
	}

	this.Output(log)
}

func (this *Logger) Debugf(format string, v ...interface{}) {
	var buf bytes.Buffer
	var log string

	message := fmt.Sprintf(format, v...)
	if this.formatCallback != nil {
		log = this.formatCallback(this, DebugLevel, message)
	} else {
		buf.WriteString("[debug] ")
		buf.WriteString(message)
		log = buf.String()
	}

	this.Output(log)
}

func (this *Logger) StartTimer(f func(*Logger)) {
	go func() {
		for {
			f(this)
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

func (this *Logger) Record(logDir string, keepDays int) {
	this.logDir = logDir
	this.keepDays = keepDays
	this.logs = make([]string, 0)
	this.StartTimer(timeout)
}

func (this *Logger) Close() {
	if this.curLogfile != nil {
		this.curLogfile.Close()
	}
}

func stdFormat(logger *Logger, level int, message string) string {
	var buf bytes.Buffer

	buf.WriteString("[")
	buf.WriteString(logger.Time(""))
	buf.WriteString("] ")

	switch level {
	case ErrorLevel:
		buf.WriteString("[error] ")
	case InfoLevel:
		buf.WriteString("[info] ")
	case DebugLevel:
		buf.WriteString("[debug] ")
	}

	file, line, _ := logger.FileAndLine(Lshortfile)
	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString("    ")
	buf.WriteString(message)
	return buf.String()
}

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func SetFormatCallback(formatCallback FormatCallback) {
	std.SetFormatCallback(formatCallback)
}

func Error(v ...interface{}) {
	std.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

func Info(v ...interface{}) {
	std.Info(v...)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

func Debug(v ...interface{}) {
	std.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func Record(logDir string, keepDays int) {
	std.Record(logDir, keepDays)
}

func StartTimer(f func(*Logger)) {
	std.StartTimer(f)
}

func Close() {
	std.Close()
}
