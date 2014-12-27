package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	Version    = "1.0.0"
	ServerName = "IPGeoBase-Go/" + Version
)

var (
	dirName        = flag.String("dir", "", "Data directory (required, write access needed)")
	geoBaseURL     = flag.String("url", "http://ipgeobase.ru/files/db/Main/geo_files.zip", "URL of IPGeoBase zip archive")
	lstAddress     = flag.String("listen", "localhost:7364", "IP address and port to listen")
	updateInterval = flag.Duration("interval", time.Hour, "Update interval")
	showHelp       = flag.Bool("help", false, "Show help")
	showVersion    = flag.Bool("version", false, "Show version number")
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

	go updateRoutine(conf)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		base := conf.VBase.Load().(*GeoBase)
		result := base.Find(ip)
		w.Header().Set("Server", ServerName)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(result)
	})

	conf.Log.Info("Starting server at " + conf.LstAddr.String())
	if err := http.ListenAndServe(conf.LstAddr.String(), nil); err != nil {
		conf.Log.Errorf("Serve error: %v", err)
		os.Exit(1)
	}
}
