package logger

import (
	"bytes"
	"os"
	"strconv"
	"testing"
	"time"
)

func Format(logger *Logger, level int, message string) string {
	var buf bytes.Buffer

	switch level {
	case ErrorLevel:
		buf.WriteString("[error] ")
	case InfoLevel:
		buf.WriteString("[info] ")
	case DebugLevel:
		buf.WriteString("[debug] ")
	}

	buf.WriteString("[")
	buf.WriteString(logger.Time(""))
	buf.WriteString("] ")
	file, line, _ := logger.FileAndLine(Lshortfile)
	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")
	buf.WriteString(message)
	return buf.String()
}

func TestCustom(t *testing.T) {
	log := New(os.Stdout, 0)
	log.SetFormatCallback(Format)
	log.Error("test error")
	log.Info("test info")
	log.Debug("test debug")
	log.Errorf("test errorf\n")
	log.Infof("test infof\n")
	log.Debugf("test debugf\n")
}

func TestStd(t *testing.T) {
	Error("test error")
	Info("test info")
	Debug("test debug")
	Errorf("test errorf\n")
	Infof("test infof\n")
	Debugf("test debugf\n")
}

func TestRecorde(t *testing.T) {
	Record("./log", 10)

	time.Sleep(time.Second * 20)
}
