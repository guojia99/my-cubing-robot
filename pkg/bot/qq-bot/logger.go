package qq_bot

import (
	"k8s.io/klog"

	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/log"
)

var logger log.Logger = &Logger{}

type Logger struct {
}

func (l Logger) Debug(v ...interface{}) {
	klog.Infof("Debug", v...)
}

func (l Logger) Info(v ...interface{}) {
	//klog.Info(v...)
}

func (l Logger) Warn(v ...interface{}) {
	klog.Warning(v...)
}

func (l Logger) Error(v ...interface{}) {
	klog.Error(v...)
}

func (l Logger) Debugf(format string, v ...interface{}) {
	klog.Infof("Debug"+format, v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	//klog.Infof(format, v...)
}

func (l Logger) Warnf(format string, v ...interface{}) {
	klog.Warningf(format, v...)
}

func (l Logger) Errorf(format string, v ...interface{}) {
	klog.Errorf(format, v...)
}

func (l Logger) Sync() error {
	return nil
}
