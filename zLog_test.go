package zLog

import (
	"testing"
)

func TestInit(t *testing.T) {

	err := Init(
		WithCommonLog("logs", "common.log", true),
		WithInfoLog("logs", "info.log", false),
		WithErrorLog("logs", "error.log", false),
		WithLogCut(1, 5, 30, true),
		WithTimeFormat("2006-01-02 15:04:05.000"))
	if err != nil {
		t.Error(err)
		return
	}
	Log.Info("hello world")
	Log.Error("hello world")
	Log.Debug("hello world")
	Log.Warn("hello world")
}

func TestDefaultLog(t *testing.T) {
	DefaultLog.Info("hello world")
	DefaultLog.Error("hello world")
	DefaultLog.Debug("hello world")
	DefaultLog.Warn("hello world")
}
