package core

import (
	"runtime"
	"log"
	"os"
)

func ServerRun(c *Config, callback func() int) int {
	LogerTrace.PrintIn("apply log tank", c.Log.tank)
}
