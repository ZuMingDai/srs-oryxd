package core

import (
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
	LoggerInfo.Println("start to parse config file", conf)
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

	return
}

func (s *Server) Run() (err error) {
	s.applyMultipleProcesses(GsConfig.Workers)

	for {
		runtime.GC()
		LoggerInfo.Println("go runtime gc every", GsConfig.Go.GcInterval, "seconds")
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
	LoggerTrace.Println("apply workers", workers, "and previousis", pv)
}

func (s *Server) applyLogger(c *Config) (err error) {
	if err = s.logger.Close(c); err != nil {
		return
	}
	LoggerInfo.Println("close logger ok")

	if err = s.logger.Open(c); err != nil {
		return
	}
	LoggerInfo.Println("open logger ok")

	return
}
