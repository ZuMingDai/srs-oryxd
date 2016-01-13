package main

import (
	"fmt"
	"runtime"
	"time"
)

type Server struct {
	logger *simpleLogger
}

func NewServer() *Server {
	svr := &Server{
		logger: &simpleLogger{},
	}
	GsConfig.Subscribe(svr)
	return svr
}

func (s *Server) Close() {
	GsConfig.Unsubscribe(s)
}

func (s *Server) ParseConfig(conf string) (err error) {
	GsTrace.Println("start to parse config file", conf)
	if err = GsConfig.Loads(conf); err != nil {
		return
	}
	return
}

func (s *Server) PrepareLogger() (err error) {
	if err = s.applyLogger(GsConfig); err != nil {
		return
	}
	return
}

func (s *Server) Initialize() (err error) {
	go ReloadWorker()

	c := GsConfig
	l := fmt.Sprintf("%v(%v/%v)", c.Log.Tank, c.Log.Level, c.Log.File)
	if !c.LogToFile() {
		l = fmt.Sprintf("%v(%v)", c.Log.Tank, c.Log.Level)
	}

	GsTrace.Println(fmt.Sprintf("init server ok,conf=%v, log=%v,workers=%v", c.conf, l, c.Workers))

	return
}

func (s *Server) Run() (err error) {
	s.applyMultipleProcesses(GsConfig.Workers)

	for {
		runtime.GC()
		GsInfo.Println("go runtime gc every", GsConfig.Go.GcInterval, "seconds")
		time.Sleep(time.Second * time.Duration(GsConfig.Go.GcInterval))
	}
	return
}

func (s *Server) OnReloadGlobal(scope int, cc, pc *Config) (err error) {
	if scope == ReloadWorkers {
		s.applyMultipleProcesses(cc.Workers)
	} else if scope == ReloadLog {
		s.applyLogger(cc)
	}
	return
}

func (s *Server) applyMultipleProcesses(workers int) {
	pv := runtime.GOMAXPROCS(workers)

	if pv != workers {
		GsTrace.Println("apply workers", workers, "and previousis", pv)
	} else {
		GsInfo.Println("apply workers", workers, "and previousis", pv)
	}

}

func (s *Server) applyLogger(c *Config) (err error) {
	if err = s.logger.Close(c); err != nil {
		return
	}
	GsInfo.Println("close logger ok")

	if err = s.logger.Open(c); err != nil {
		return
	}
	GsInfo.Println("open logger ok")

	return
}
