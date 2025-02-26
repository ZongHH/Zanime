package test

import (
	"fmt"
	"gateService/pkg/monitor"
	"testing"
)

func Test_Monitor(t *testing.T) {
	monitor.Init(monitor.NewLogConfig())

	for i := 0; i < 100; i++ {
		fmt.Println(i)
		monitor.Info("this is a test")
		monitor.Error("this is a test")
		monitor.Warning("this is a test")
		monitor.Info("this is a test")
	}

	monitor.Close()
}
