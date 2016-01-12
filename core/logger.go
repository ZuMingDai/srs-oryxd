package core

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	logLabel      = "[gsrs]"
	logInfoLabel  = logLabel + "[info]"
	logTraceLabel = logLabel + "[trace]"
	logWarnLabel  = logLabel + "[warn]"
	logErrorLabel = logLabel + "[error]"
)

//the application loggers
//info, the verbose info level, very detail log, the lowest lever, to discard
var LoggerInfo Logger = log.New(ioutil.Discard, logInfoLabel, log.LstdFlags)

//
var LoggerTrace Logger = log.New(os.Stdout, logTraceLabel, log.LstdFlags)

//
var LoggerWarn Logger = log.New(os.Stderr, logWarnLabel, log.LstdFlags)

//
var LoggerError Logger = log.New(os.Stderr, logErrorLabel, log.LstdFlags)

type Logger interface {
	Print(a ...interface{})
	Println(a ...interface{})
	Printf(format string, a ...interface{})
}

type simpleLogger struct {
	file *os.File
}

func (l *simpleLogger) Open(c *Config) (err error) {
	LoggerInfo.Println("apply log tank", c.Log.Tank)
	LoggerInfo.Println("apply log level", c.Log.Level)

	if c.LogToFile() {
		LoggerTrace.Println("apply log", c.Log.Tank, c.Log.Level, c.Log.File)
		LoggerTrace.Println("please see detail of log:tailf", c.Log.File)

		if l.file, err = os.OpenFile(c.Log.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
			LoggerError.Println("open log file", c.Log.File, "failed,err is", err)
			return
		} else {
			LoggerInfo = log.New(c.LogTank("info", l.file), logInfoLabel, log.LstdFlags)
			LoggerTrace = log.New(c.LogTank("trace", l.file), logTraceLabel, log.LstdFlags)
			LoggerWarn = log.New(c.LogTank("warn", l.file), logWarnLabel, log.LstdFlags)
			LoggerError = log.New(c.LogTank("error", l.file), logErrorLabel, log.LstdFlags)
		}

	} else {
		LoggerTrace.Println("apply log", c.Log.Tank, c.Log.Level)

		LoggerInfo = log.New(c.LogTank("info", os.Stdout), logInfoLabel, log.LstdFlags)
		LoggerTrace = log.New(c.LogTank("trace", os.Stdout), logTraceLabel, log.LstdFlags)
		LoggerWarn = log.New(c.LogTank("warn", os.Stderr), logWarnLabel, log.LstdFlags)
		LoggerError = log.New(c.LogTank("error", os.Stderr), logErrorLabel, log.LstdFlags)
	}
	return
}

func (l *simpleLogger) Close(c *Config) (err error) {
	if l.file == nil {
		return
	}

	LoggerWarn = log.New(os.Stderr, logWarnLabel, log.LstdFlags)

	if err = l.file.Close(); err != nil {
		LoggerWarn.Println("gracefully close log file", c.Log.File, "failed,err is", err)

	} else {
		LoggerWarn.Println("close log file", c.Log.File, "ok")
	}

	return
}
