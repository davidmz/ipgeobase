package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/davidmz/ipgeobase/mmc"
)

const (
	Version    = "1.1.1"
	ServerName = "IPGeoBase-Go/" + Version
)

var (
	dirName        = flag.String("dir", "", "Data directory (required, write access needed unless passive mode)")
	geoBaseURL     = flag.String("url", "http://ipgeobase.ru/files/db/Main/geo_files.zip", "URL of IPGeoBase zip archive")
	lstAddress     = flag.String("listen", "localhost:7364", "IP address and port to listen")
	updateInterval = flag.Duration("interval", time.Hour, "Update interval")
	logFile        = flag.String("log", "", "File name for log output")
	passiveMode    = flag.Bool("passive", false, "Passive mode: do not write to data directory")
	debugLevel     = flag.Bool("debug", false, "Debug level log")
	showHelp       = flag.Bool("help", false, "Show help")
	showVersion    = flag.Bool("version", false, "Show version number")
	asMemcache     = flag.Bool("memcache", false, "Serve memcache protocol")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	conf := GetConfig()

	fileExists, err := FileExists(conf.Dir.baseFile())
	if err != nil {
		conf.Log.Errorf("File access error: %v", err)
		os.Exit(1)
	}

	if !fileExists {
		if conf.PassiveMode {
			conf.Log.Error("Database file not exists")
			os.Exit(1)
		}

		conf.Log.Info("IPGeoBase file not exists, downloading...")
		for {
			os.Remove(conf.Dir.eTagFile())
			if _, err := downloadGeoBase(conf); err != nil {
				conf.Log.Warnf("Download error: %v", err)
				time.Sleep(5 * time.Second)
			} else {
				conf.Log.Info("File was downloaded")
				break
			}
		}
	}

	if err := parseAndStoreGeoBase(conf); err != nil {
		conf.Log.Errorf("Parse error: %v", err)
		os.Exit(1)
	}

	if conf.PassiveMode {
		go passiveUpdateRoutine(conf)
	} else {
		go updateRoutine(conf)
	}

	conf.Log.Info("Starting TCP server at " + conf.LstAddr.String())

	ln, err := net.Listen("tcp", conf.LstAddr.String())
	if err != nil {
		conf.Log.Errorf("Server error: %v", err)
		os.Exit(1)
	}

	if *asMemcache {
		h := &MemcacheHandler{conf}
		for {
			conn, err := ln.Accept()
			conf.Log.Debugf("Memcache connect from %q", conn.RemoteAddr())
			if err != nil {
				conf.Log.Errorf("Serve error: %v", err)
				os.Exit(1)
			}
			go mmc.NewSession(conn, h)
		}

	} else {

		s := &http.Server{
			Handler:        &HttpHandler{conf},
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 8 << 10,
		}

		if err := s.Serve(ln); err != nil {
			conf.Log.Errorf("HTTP serve error: %v", err)
			os.Exit(1)
		}

	}
}
