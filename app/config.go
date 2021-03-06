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
package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZuMingDai/srs-oryxd/core"
)

const (
	ReloadWorkers = iota
	//	ReloadListen
	ReloadLog
)

type ReloadHandler interface {
	OnReloadGlobal(scope int, cc, pc *Config) error
}

type Config struct {
	Workers int `json:"workers"`
	Listen  int `json:"listen"`

	Go struct {
		GcInterval int `json:"gc_interval"`
	}

	Log struct {
		Tank  string `json:"tank"`
		Level string `json:"level"`
		File  string `json:"file"`
	} `json:"log"`

	conf           string          `json:"-"`
	reloadHandlers []ReloadHandler `json:"-"`
}

var GsConfig = NewConfig()

func NewConfig() *Config {
	c := &Config{
		reloadHandlers: []ReloadHandler{},
	}
	//default to use 1 cpu
	c.Workers = core.Workers
	c.Listen = core.RtmpListen
	c.Go.GcInterval = core.GcIntervalSeconds

	c.Log.Tank = "file"
	c.Log.Level = "trace"
	c.Log.File = "gsrs.log"

	return c
}

func (c *Config) Loads(conf string) error {
	c.conf = conf

	var s []byte
	if f, err := os.Open(conf); err != nil {
		return err
	} else if s, err = ioutil.ReadAll(f); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(s), c); err != nil {
		return err
	}
	return c.Validate()
}

func (c *Config) Validate() error {
	if c.Log.Level == "info" {
		core.GsWarn.Println("info level hurts performance")
	}
	if c.Workers <= 0 || c.Workers > 64 {
		return errors.New(fmt.Sprintf("workers must in(0,64], actual is %v", c.Workers))
	}
	if c.Listen <= 0 || c.Listen > 65535 {
		return errors.New(fmt.Sprintf("listen must in(0,65535], actual is %v", c.Listen))
	}
	if c.Go.GcInterval <= 0 || c.Go.GcInterval > 24*3600 {
		return errors.New(fmt.Sprintf("go gc_interval must in (0, 24*3600], actual is %v", c.Go.GcInterval))
	}
	if c.Log.Level != "info" && c.Log.Level != "trace" && c.Log.Level != "warn" && c.Log.Level != "error" {
		return errors.New(fmt.Sprintf("log.level must be info/trace/warn/error, actual is %v", c.Log.Level))
	}
	if c.Log.Tank != "console" && c.Log.Tank != "file" {
		return errors.New(fmt.Sprintf("log.tank must be console/file, actual is %v", c.Log.Tank))
	}
	if c.Log.Tank == "file" && len(c.Log.File) == 0 {
		return errors.New("log.file must not be empty for file tank")
	}
	return nil

}

/*
func (c *Config) Json() (string, error) {
	if b, err := json.Marshal(c); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
*/
func (c *Config) LogToFile() bool {
	return c.Log.Tank == "file"
}

func (c *Config) LogTank(level string, dw io.Writer) io.Writer {
	if c.Log.Level == "info" {
		return dw
	}
	if c.Log.Level == "trace" {
		if level == "info" {
			return ioutil.Discard
		}
		return dw
	}
	if c.Log.Level == "warn" {
		if level == "info" || level == "trace" {
			return ioutil.Discard
		}
		return dw
	}
	if c.Log.Level == "error" {
		if level != "error" {
			return ioutil.Discard
		}
		return dw
	}

	return ioutil.Discard
}

func (c *Config) Subscribe(h ReloadHandler) {
	for _, v := range c.reloadHandlers {
		if v == h {
			return
		}
	}
	c.reloadHandlers = append(c.reloadHandlers, h)
}

func (c *Config) Unsubscribe(h ReloadHandler) {
	for i, v := range c.reloadHandlers {
		if v == h {
			c.reloadHandlers = append(c.reloadHandlers[:i], c.reloadHandlers[i+1:]...)
			return
		}
	}
}

func (pc *Config) Reload(cc *Config) (err error) {
	if cc.Workers != pc.Workers {
		for _, h := range cc.reloadHandlers {
			if err = h.OnReloadGlobal(ReloadWorkers, cc, pc); err != nil {
				return
			}
		}
		core.GsTrace.Println("reload apply workers ok")
	} else {
		core.GsInfo.Println("reload ignore workers")
	}

	if cc.Log.File != pc.Log.File || cc.Log.Level != pc.Log.Level || cc.Log.Tank != pc.Log.Tank {
		for _, h := range cc.reloadHandlers {
			if err = h.OnReloadGlobal(ReloadLog, cc, pc); err != nil {
				return
			}
		}
		core.GsTrace.Println("reload apply log ok")
	} else {
		core.GsInfo.Println("reload ignore log")
	}

	return
}
func (c *Config) reloadCycle(wc WorkerContainer) {
	signals := make(chan os.Signal, 1)
	//1:SIGHUP
	signal.Notify(signals, syscall.Signal(1))

	core.GsTrace.Println("wait for reload signals:kill -1", os.Getpid())
	for {
		select {
		case signal := <-signals:
			core.GsTrace.Println("start reload by", signal)

			if err := c.doReload(); err != nil {
				core.GsError.Println("quit for reload failed. err is", err)

				wc.Quit()

				return
			}
		case <-wc.QC():
			core.GsWarn.Println("user stop reload.")

			wc.Quit()

			return
		}
	}

}

func (c *Config) doReload() (err error) {

	pc := c
	cc := NewConfig()
	cc.reloadHandlers = pc.reloadHandlers[:]
	if err = cc.Loads(c.conf); err != nil {
		core.GsError.Println("reload config failed,err is", err)
		return
	}
	core.GsInfo.Println("reload parse fresh config ok")
	if err = pc.Reload(cc); err != nil {
		core.GsError.Println("apply reload failed,err is", err)
		return
	}
	core.GsInfo.Println("reload completed work")

	GsConfig = cc
	core.GsTrace.Println("reload config ok")

	return
}
