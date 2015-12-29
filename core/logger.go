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
var LoggerInfo Logger = log.New(ioutil.Discard, logLabel, log.LstdFlags)

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
