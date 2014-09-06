package log

import "github.com/golang/glog"

type GlogRecorder struct {
	Recorder
}

func (l *GlogRecorder) Info(m string) {
	glog.Info(m)
}

func (l *GlogRecorder) Debug(m string) {
	if glog.V(4) {
		glog.Info(m)
	}
}

func (l *GlogRecorder) Error(m string) {
	glog.Error(m)
}
func (l *GlogRecorder) Warn(m string) {
	glog.Warning(m)
}
