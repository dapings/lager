package glog

import (
	"testing"
)

func TestAll(t *testing.T) {
	Errorf("%s %s", "foo", "bar")
	Error("foo", "bar")
	Errorln("foo", "bar")
	ErrorDepth(2, "foo", "bar")
	
	Warningf("%s %s", "foo", "bar")
	Warning("foo", "bar")
	Warningln("foo", "bar")
	WarningDepth(2, "foo", "bar")
	
	Infof("%s %s", "foo", "bar")
	Info("foo", "bar")
	Infoln("foo", "bar")
	InfoDepth(2, "foo", "bar")
	
	for i := 0; i < 10; i++ {
		V(Level(i)).Infof("%s %s", "foo", "bar")
		V(Level(i)).Info("foo", "bar")
		V(Level(i)).Infoln("foo", "bar")
	}
	
	Flush()
}
