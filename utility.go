package main

/*
import (
	"log"
	"os"
	"runtime"
)

func ServerRun(c *Config, callback func() int) int {
	LoggerTrace.Println("apply log tank", c.Log.Tank)
	LoggerTrace.Println("apply log level", c.Log.Level)
	if c.LogToFile() {
		LoggerTrace.Println("apply log file ", c.Log.File)
		if f, err := os.OpenFile(c.Log.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
			LoggerError.Println("open log file", c.Log.File, "failed, err is", err)
			return -1
		} else {
			defer func() {
				err = f.Close()
				LoggerWarn = log.New(os.Stderr, "[gsrs][warn]", log.LstdFlags)
				if err != nil {
					LoggerWarn.Println("gracefully close log file", c.Log.File, "failed, err is", err)
				} else {
					LoggerWarn.Println("close log file", c.Log.File, "ok")
				}
			}()
			LoggerInfo = log.New(c.LogTank("info", f), logInfoLabel, log.LstdFlags)
			LoggerTrace = log.New(c.LogTank("trace", f), logTraceLabel, log.LstdFlags)
			LoggerWarn = log.New(c.LogTank("warn", f), logWarnLabel, log.LstdFlags)
			LoggerError = log.New(c.LogTank("error", f), logErrorLabel, log.LstdFlags)
		}
		LoggerTrace.Println("please see detail of log: tailf", c.Log.File)
	} else {
		LoggerInfo = log.New(c.LogTank("info", os.Stdout), logInfoLabel, log.LstdFlags)
		LoggerTrace = log.New(c.LogTank("trace", os.Stdout), logTraceLabel, log.LstdFlags)
		LoggerWarn = log.New(c.LogTank("warn", os.Stderr), logWarnLabel, log.LstdFlags)
		LoggerError = log.New(c.LogTank("error", os.Stderr), logErrorLabel, log.LstdFlags)
	}
	LoggerTrace.Println("apply workers", c.Workers)
	runtime.GOMAXPROCS(c.Workers)

	return callback()
}
*/
