package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ZuMingDai/srs-oryxd/app"
	"github.com/ZuMingDai/srs-oryxd/core"
)

//the startup argv
//	-c conf/srs.json
//	--c conf/srs.json
//	-c=conf/srs.json
//	--c=conf/srs.json
var confFile = flag.String("c", "D:\\mygo\\src\\github.com\\ZuMingDai\\srs-oryxd\\conf\\srs.json", "the config file.")

func run() int {

	//	core.LoggerTrace.Println(fmt.Sprintf("GO-SRS/%v is a golang implementation os SRS", core.Version))
	flag.Parse()

	svr := app.NewServer()
	defer svr.Close()

	if err := svr.ParseConfig(*confFile); err != nil {
		core.GsError.Println("parse config from", *confFile, "failed,err is", err)
		return -1
	}

	if err := svr.PrepareLogger(); err != nil {
		core.GsError.Println("prepare logger failed,err is", err)
		return -1
	}

	/*
		core.LoggerInfo.Println("star to parse config file", confFile)
		if err := core.GsConfig.Loads(confFile); err != nil {
			core.LoggerError.Println("parse config", confFile, "failed,err is", err)
			return -1
		}
		go core.GsConfig.ReloadWorker(confFile)
	*/
	core.GsTrace.Println("Copyright (c) 2013-2015 SRS(simple-rtmp-server)")

	core.GsTrace.Println(fmt.Sprintf("GO-SRS/%v is a golang implementation of SRS.", core.Version()))
	if err := svr.Initialize(); err != nil {
		core.GsError.Println("initialize server failed,err is", err)
		return -1
	}

	if err := svr.Run(); err != nil {
		core.GsError.Println("run server failed, err is", err)
		return -1
	}

	return 0
}

func main() {
	ret := run()
	os.Exit(ret)
}
