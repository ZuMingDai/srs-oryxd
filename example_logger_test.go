package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestBaseceLogger(t *testing.T) {

	GsInfo = log.New(ioutil.Discard, logLabel, log.LstdFlags)
	GsTrace = log.New(ioutil.Discard, logLabel, log.LstdFlags)
	GsWarn = log.New(ioutil.Discard, logLabel, log.LstdFlags)
	GsError = log.New(ioutil.Discard, logLabel, log.LstdFlags)

	GsInfo.Println("test logger info ok.")
	GsTrace.Println("test logger trace ok.")
	GsWarn.Println("test logger warn ok.")
	GsError.Println("test logger error ok.")
}
