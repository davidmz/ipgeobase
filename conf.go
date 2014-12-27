package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/ivpusic/golog"
)

type Config struct {
	Log   *golog.Logger
	VBase *atomic.Value

	Dir         *DataDir
	GeoBaseUrl  string
	LstAddr     net.Addr
	UpdInterval time.Duration
}

func ExitIfError(conf *Config, err error) {
	if err != nil {
		conf.Log.Error(err.Error())
		flag.Usage()
		os.Exit(1)
	}
}

func GetConfig() *Config {
	conf := &Config{
		Log:   golog.GetLogger("ipgeo"),
		VBase: &atomic.Value{},
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	// dirName

	if *dirName == "" {
		ExitIfError(conf, fmt.Errorf("Directory name is required"))
	}

	// Write test
	probeFileName := filepath.Join(*dirName, "probe")
	if f, err := os.Create(probeFileName); err != nil {
		ExitIfError(conf, err)
	} else {
		f.Close()
		os.Remove(probeFileName)
	}

	conf.Dir = &DataDir{*dirName}

	// geoBaseURL

	u, err := url.Parse(*geoBaseURL)
	ExitIfError(conf, err)

	conf.GeoBaseUrl = u.String()

	// lstAddress

	a, err := net.ResolveTCPAddr("tcp", *lstAddress)
	ExitIfError(conf, err)

	conf.LstAddr = a

	// updateInterval

	if *updateInterval <= 0 {
		ExitIfError(conf, fmt.Errorf("Invalid update interval"))
	}

	conf.UpdInterval = *updateInterval

	return conf
}
