package log

import "fmt"

//GlobalRecorder is pointer to the global log recorder
var GlobalRecorder Recorder = &GlogRecorder{}

//Recorder is the recorder struct containing the associated functions
type Recorder interface {
	Info(m string)
	Debug(m string)
	Warn(m string)
	Error(m string)
	Flush()
}

//Debug records the supplied debug string
func Debug(m string) {
	GlobalRecorder.Debug(m)
}

//Info records the supplied Info string
func Info(m string) {
	GlobalRecorder.Info(m)
}

//Warn records the supplied Warn string
func Warn(m string) {
	GlobalRecorder.Warn(m)
}

//Error records the supplied Error string
func Error(m string) {
	GlobalRecorder.Error(m)
}

//Debugf records the supplied Debug string including any arguments (i.e. Printf)
func Debugf(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Debug(mf)
}

//Infof records the supplied Info string including any arguments (i.e. Printf)
func Infof(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Info(mf)
}

//Warnf records the supplied Warn string including any arguments (i.e. Printf)
func Warnf(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Warn(mf)
}

//Errorf records the supplied Error string including any arguments (i.e. Printf)
func Errorf(m string, args ...interface{}) {
	mf := fmt.Sprintf(m, args...)
	GlobalRecorder.Error(mf)
}

//Init initializes the recorder
func Init(recorder Recorder) {
	GlobalRecorder = recorder
}

//Flush flushes the recorder
func Flush() {
	GlobalRecorder.Flush()
}
