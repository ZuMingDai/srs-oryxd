/*
The MIT License (MIT)

Copyright (c) 2013-2015 SRS(simple-rtmp-server)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package core

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	logLabel      = "[gsrs]"
	LogInfoLabel  = logLabel + "[info]"
	LogTraceLabel = logLabel + "[trace]"
	LogWarnLabel  = logLabel + "[warn]"
	LogErrorLabel = logLabel + "[error]"
)

//the application loggers
//info, the verbose info level, very detail log, the lowest lever, to discard
var GsInfo Logger = log.New(ioutil.Discard, LogInfoLabel, log.LstdFlags)

//
var GsTrace Logger = log.New(os.Stdout, LogTraceLabel, log.LstdFlags)

//
var GsWarn Logger = log.New(os.Stderr, LogWarnLabel, log.LstdFlags)

//
var GsError Logger = log.New(os.Stderr, LogErrorLabel, log.LstdFlags)

type Logger interface {
	//	Print(a ...interface{})
	Println(a ...interface{})
	//Printf(format string, a ...interface{})
}

/*
type simpleLogger struct {
	file *os.File
}

func (l *simpleLogger) Open(c *Config) (err error) {
	GsInfo.Println("apply log tank", c.Log.Tank)
	GsInfo.Println("apply log level", c.Log.Level)

	if c.LogToFile() {
		GsTrace.Println("apply log", c.Log.Tank, c.Log.Level, c.Log.File)
		GsTrace.Println("please see detail of log:tailf", c.Log.File)

		if l.file, err = os.OpenFile(c.Log.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
			GsError.Println("open log file", c.Log.File, "failed,err is", err)
			return
		} else {
			GsInfo = log.New(c.LogTank("info", l.file), logInfoLabel, log.LstdFlags)
			GsTrace = log.New(c.LogTank("trace", l.file), logTraceLabel, log.LstdFlags)
			GsWarn = log.New(c.LogTank("warn", l.file), logWarnLabel, log.LstdFlags)
			GsError = log.New(c.LogTank("error", l.file), logErrorLabel, log.LstdFlags)
		}

	} else {
		GsTrace.Println("apply log", c.Log.Tank, c.Log.Level)

		GsInfo = log.New(c.LogTank("info", os.Stdout), logInfoLabel, log.LstdFlags)
		GsTrace = log.New(c.LogTank("trace", os.Stdout), logTraceLabel, log.LstdFlags)
		GsWarn = log.New(c.LogTank("warn", os.Stderr), logWarnLabel, log.LstdFlags)
		GsError = log.New(c.LogTank("error", os.Stderr), logErrorLabel, log.LstdFlags)
	}
	return
}

func (l *simpleLogger) Close(c *Config) (err error) {
	if l.file == nil {
		return
	}

	GsWarn = log.New(os.Stderr, logWarnLabel, log.LstdFlags)

	if err = l.file.Close(); err != nil {
		GsWarn.Println("gracefully close log file", c.Log.File, "failed,err is", err)

	} else {
		GsWarn.Println("close log file", c.Log.File, "ok")
	}
	return
}
*/
