package app

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/ZuMingDai/srs-oryxd/core"
)

//the container for all worker
//which provides the quit and cleaup methods
type WorkerContainer interface {
	//get the quit channel
	//worker can fetch the quit signal
	//or push a quit signal to channel
	QC() chan bool
	//notify the container to quit
	Quit()
	//fork a new goroutine with work container
	//the param can be a globle func or object method
	GFork(func(WorkerContainer))
}

type Server struct {
	//signal handler
	sigs chan os.Signal
	//whether closed
	closed  bool
	closing chan bool
	//for system internal to notify quit
	quit chan bool
	wg   sync.WaitGroup
	//logger
	logger *simpleLogger
	//the locker for state,for instance,the closed.
	lock sync.Mutex
}

func NewServer() *Server {
	svr := &Server{
		sigs:    make(chan os.Signal, 1),
		closed:  true,
		closing: make(chan bool, 1),
		quit:    make(chan bool, 1),
		logger:  &simpleLogger{},
	}
	GsConfig.Subscribe(svr)
	return svr
}

//notify server to stop and wait for cleanup
//TODO:FIXME:should return a chan to support sync timeout close.
func (s *Server) Close() {

	//wait for stopped
	s.lock.Lock()
	defer s.lock.Unlock()

	//close
	if s.closed {
		core.GsInfo.Println("server already colsed.")
		return
	}

	//notify to close
	core.GsInfo.Println("notify server to stop.")
	select {
	case s.quit <- true:
	default:
	}

	//wait for closed.
	<-s.closing

	//do cleanup when stopped.
	GsConfig.Unsubscribe(s)

	//ok,closed
	s.closed = true
	core.GsInfo.Println("server closed.")

}

func (s *Server) ParseConfig(conf string) (err error) {
	core.GsTrace.Println("start to parse config file", conf)
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
	//install signals
	//TODO:FIXME: when process the current signal,others may drop
	signal.Notify(s.sigs)

	//reload goroutine
	s.GFork(GsConfig.reloadCycle)

	c := GsConfig
	l := fmt.Sprintf("%v(%v/%v)", c.Log.Tank, c.Log.Level, c.Log.File)
	if !c.LogToFile() {
		l = fmt.Sprintf("%v(%v)", c.Log.Tank, c.Log.Level)
	}

	core.GsTrace.Println(fmt.Sprintf("init server ok,conf=%v, log=%v,workers=%v", c.conf, l, c.Workers))

	return
}

func (s *Server) Run() (err error) {
	func() {
		//when running,the state cannot changed.
		s.lock.Lock()
		defer s.lock.Unlock()

		//set to running
		s.closed = false
	}()

	// when terminated,notify the chan
	defer func() {
		select {
		case s.closing <- true:
		default:
		}
	}()

	core.GsInfo.Println("server runing")

	//run server,apply setting
	s.applyMultipleProcesses(GsConfig.Workers)

	for {
		select {
		case signal := <-s.sigs:
			core.GsTrace.Println("got signal", signal)
			switch signal {
			case os.Interrupt:
				//SIGINT
				s.Quit()
			case syscall.SIGTERM:
				//SIGTERM
				s.Quit()
			}
		case <-s.QC():
			s.Quit()
			s.wg.Wait()
			core.GsWarn.Println("server quit.")
			return
		case <-time.After(time.Second * time.Duration(GsConfig.Go.GcInterval)):
			runtime.GC()
			core.GsInfo.Println("go runtime gc every", GsConfig.Go.GcInterval, "seconds")
		}

	}
	return
}

//interface WorkContainer
func (s *Server) QC() chan bool {
	return s.quit
}

func (s *Server) Quit() {
	select {
	case s.quit <- true:
	default:
	}
}

func (s *Server) GFork(f func(WorkerContainer)) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		defer func() {
			if r := recover(); r != nil {
				core.GsError.Println("woker panic:", r)
				s.Quit()
			}
		}()
		f(s)
		core.GsTrace.Println("worker terminated.")
	}()
}

//interface ReloadGlobaler
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
		core.GsTrace.Println("apply workers", workers, "and previousis", pv)
	} else {
		core.GsInfo.Println("apply workers", workers, "and previousis", pv)
	}

}

func (s *Server) applyLogger(c *Config) (err error) {
	if err = s.logger.close(c); err != nil {
		return
	}
	core.GsInfo.Println("close logger ok")

	if err = s.logger.open(c); err != nil {
		return
	}
	core.GsInfo.Println("open logger ok")

	return
}
