package main

import (
	"fmt"
	"testing"
)

func TestBasicConfig(t *testing.T) {
	c := NewConfig()

	if c.Workers != Workers {
		t.Error("workers failed.")
	}

	if c.Listen != RtmpListen {
		t.Error("listen failed.")
	}

	if c.Go.GcInterval != GcIntervalSeconds {
		t.Error("go gc interval failed.")
	}

	if c.Log.Tank != "file" {
		t.Error("log tank failed.")
	}

	if c.Log.Level != "trace" {
		t.Error("log level failed.")
	}

	if c.Log.File != "gsrs.log" {
		t.Error("log file failed.")
	}
}

func BenchmarkConfigBasic(b *testing.B) {
	pc := NewConfig()
	cc := NewConfig()
	if err := pc.Reload(cc); err != nil {
		b.Error("reload failed.")
	}
}

func ExampleConfig_Loads() {
	c := NewConfig()
	fmt.Println("listen at", c.Listen)
}
