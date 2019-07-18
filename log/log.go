package logger

import (
	"fmt"
	"time"

	scrlog "github.com/Myriad-Dreamin/screenrus"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New()

type NilWriter struct {
}

func (nw *NilWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

var predo = func(log.Level) bool { return false }
var (
	FlushTime = time.Millisecond * 100
)

func ResetLine() {
	fmt.Printf("\033[2K\r\r")
}
func ResetToLastLine() {
	fmt.Printf("\033[A\r\033[2K\r")
}

func ToNextLine() {
	fmt.Printf("\033[E\r")
}

func init() {
	var screenLog, _ = scrlog.NewScreenLogPlugin(nil)
	Logger.SetOutput(new(NilWriter))
	Logger.AddHook(screenLog)
	var loggerTime = time.Now()
	var ninf = true
	predo = func(l log.Level) bool {
		if l == log.InfoLevel {
			if time.Now().Sub(loggerTime) > FlushTime {
				if ninf {
					ninf = false
				} else {
					ResetToLastLine()
				}
				loggerTime = time.Now()
				return false
			} else {
				return true
			}
		} else {
			ninf = true
			return false
		}
	}
}

func Logf(level log.Level, format string, args ...interface{}) {
	if predo(level) {
		return
	}
	if Logger.IsLevelEnabled(level) {
		Logger.Logf(level, format, args...)
	}
}

func Logln(level log.Level, args ...interface{}) {
	if predo(level) {
		return
	}
	if Logger.IsLevelEnabled(level) {
		Logger.Logln(level, args...)
	}
}

func Log(level log.Level, args ...interface{}) {
	if predo(level) {
		return
	}
	if Logger.IsLevelEnabled(level) {
		Logger.Log(level, args...)
	}
}

func Printf(format string, args ...interface{}) {
	if predo(2) {
		return
	}
	Logger.Printf(format, args...)
}

func Println(args ...interface{}) {
	if predo(2) {
		return
	}
	Logger.Println(args...)
}

func Warningln(args ...interface{}) {
	if predo(2) {
		return
	}
	Logger.Warnln(args...)
}

func Tracef(format string, args ...interface{}) {
	Logf(log.TraceLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	Logf(log.DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	Logf(log.InfoLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	Logf(log.WarnLevel, format, args...)
}

func Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logf(log.ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Logf(log.FatalLevel, format, args...)
	Logger.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	Logf(log.PanicLevel, format, args...)
}

func Trace(args ...interface{}) {
	Log(log.TraceLevel, args...)
}

func Debug(args ...interface{}) {
	Log(log.DebugLevel, args...)
}

func Info(args ...interface{}) {
	Log(log.InfoLevel, args...)
}

func Warn(args ...interface{}) {
	Log(log.WarnLevel, args...)
}

func Warning(args ...interface{}) {
	Warn(args...)
}

func Error(args ...interface{}) {
	Log(log.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	Log(log.FatalLevel, args...)
	Logger.Exit(1)
}

func Panic(args ...interface{}) {
	Log(log.PanicLevel, args...)
}

func Traceln(args ...interface{}) {
	Logln(log.TraceLevel, args...)
}

func Debugln(args ...interface{}) {
	Logln(log.DebugLevel, args...)
}

func Infoln(args ...interface{}) {
	Logln(log.InfoLevel, args...)
}

func Warnln(args ...interface{}) {
	Logln(log.WarnLevel, args...)
}

func Errorln(args ...interface{}) {
	Logln(log.ErrorLevel, args...)
}

func Fatalln(args ...interface{}) {
	Logln(log.FatalLevel, args...)
	Logger.Exit(1)
}

func Panicln(args ...interface{}) {
	Logln(log.PanicLevel, args...)
}
