package log

import "fmt"

var GlobalRecorder Recorder = &GlogRecorder{}

type Recorder interface {
	Info(m string)
	Debug(m string)
	Warn(m string)
	Error(m string)
}

func Debug(m string) {
	GlobalRecorder.Debug(m)
}

func Info(m string) {
	GlobalRecorder.Info(m)
}

func Warn(m string) {
	GlobalRecorder.Warn(m)
}

func Error(m string) {
	GlobalRecorder.Error(m)
}

func Debugf(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Debug(mf)
}

func Infof(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Info(mf)
}

func Warnf(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Warn(mf)
}

func Errorf(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Error(mf)
}

func Init(recorder Recorder) {
	GlobalRecorder = recorder
}
