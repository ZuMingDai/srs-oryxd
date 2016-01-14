package main

import (
	"log"
	"strings"
	"testing"
)

// convert a func to interface io.Writer
type WriterFunc func(p []byte) (n int, err error)

// for io.Writer
func (f WriterFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

func TestBasicLogger(t *testing.T) {
	var tank string
	var writer = func(p []byte) (n int, err error) {
		tank = string(p)
		return len(tank), nil
	}

	GsInfo = log.New(WriterFunc(writer), logInfoLabel, log.LstdFlags)
	GsTrace = log.New(WriterFunc(writer), logTraceLabel, log.LstdFlags)
	GsWarn = log.New(WriterFunc(writer), logWarnLabel, log.LstdFlags)
	GsError = log.New(WriterFunc(writer), logErrorLabel, log.LstdFlags)

	GsInfo.Println("test logger.")
	if !strings.HasPrefix(tank, "[gsrs][info]") {
		t.Error("logger format failed.")
	}
	if !strings.HasSuffix(tank, "test logger.\n") {
		t.Error("logger format failed. tank is", tank)
	}

	GsTrace.Println("test logger.")
	if !strings.HasPrefix(tank, "[gsrs][trace]") {
		t.Error("logger format failed.")
	}
	if !strings.HasSuffix(tank, "test logger.\n") {
		t.Error("logger format failed. tank is", tank)
	}

	GsWarn.Println("test logger.")
	if !strings.HasPrefix(tank, "[gsrs][warn]") {
		t.Error("logger format failed.")
	}
	if !strings.HasSuffix(tank, "test logger.\n") {
		t.Error("logger format failed. tank is", tank)
	}

	GsError.Println("test logger.")
	if !strings.HasPrefix(tank, "[gsrs][error]") {
		t.Error("logger format failed.")
	}
	if !strings.HasSuffix(tank, "test logger.\n") {
		t.Error("logger format failed. tank is", tank)
	}

	//	t.Error("logger format error.")
}
