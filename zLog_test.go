package zLog

import (
	"testing"
)

func TestInit(t *testing.T) {

	err := Init(WithCommonLog("logs", "common.log", true))
	if err != nil {
		t.Error(err)
		return
	}
	Log.Info("hello world")
	Log.Error("hello world")
	Log.Debug("hello world")
	Log.Warn("hello world")
}
