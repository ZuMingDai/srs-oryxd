package core

import (
	"log"
	"os"
	"runtime"
)

func ServerRun(c *Config, callback func() int) int {
	LogerTrace.PrintIn("apply log tank", c.Log.Tank)
}
