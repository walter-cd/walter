package log

import "github.com/ainoya/glog"

type GlogRecorder struct {
	Recorder
}

func (l *GlogRecorder) Debug(m string) {
	glog.Debug(m)
}

func (l *GlogRecorder) Info(m string) {
	glog.Info(m)
}

func (l *GlogRecorder) Error(m string) {
	glog.Error(m)
}

func (l *GlogRecorder) Warn(m string) {
	glog.Warning(m)
}

func (l *GlogRecorder) Flush() {
	glog.Flush()
}
