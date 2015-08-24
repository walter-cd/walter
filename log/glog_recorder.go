package log

import "github.com/ainoya/glog"

//GlogRecorder struct contains the recorder
type GlogRecorder struct {
	Recorder
}

//Debug records the supplied debug string
func (l *GlogRecorder) Debug(m string) {
	glog.Debug(m)
}

//Info records the supplied Info string
func (l *GlogRecorder) Info(m string) {
	glog.Info(m)
}

//Error records the supplied error string
func (l *GlogRecorder) Error(m string) {
	glog.Error(m)
}

//Warn records the supplied warn string
func (l *GlogRecorder) Warn(m string) {
	glog.Warning(m)
}

//Flush flushes the recorder
func (l *GlogRecorder) Flush() {
	glog.Flush()
}
