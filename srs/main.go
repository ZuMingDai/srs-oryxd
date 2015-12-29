package main

import (
	"flag"
	"fmt"
	"go-oryxd/core"
	"os"
)

//the startup argv
//	-c conf/srs.json
//	--c conf/srs.json
//	-c=conf/srs.json
//	--c=conf/srs.json
var conffile = *flag.String("c", "conf/srs.json", "the config file.")

func run() int {
	core.LoggerTrace.Println(fmt.Sprintf("GO-SRS/%v is a golang implementation os SRS", core.Version))
	flag.Parse()

	conf := &core.Config{}
	core.LoggerInfo.Println("start to parse config file", conffile)

	if err := conf.Loads(conffile); err != nil {
		core.LoggerError.Println("parse config", conffile, "failed, err is", err)
		return -1
	}

	core.LoggerTrace.Println("Copyright (c) 2013-2015 SRS(simple-rtmp-server")
	return core.ServerRun(conf, func() int {
		return 0
	})
}

func main() {
	ret := run()
	os.Exit(ret)
}
