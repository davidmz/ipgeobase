package main

import (
	"os"
	"time"
)

func updateRoutine(conf *Config) {
	clock := time.Tick(conf.UpdInterval)
	for _ = range clock {
		conf.Log.Info("Looking for new base...")

		wasChanged, err := downloadGeoBase(conf)
		if err != nil {
			conf.Log.Warnf("Download error: %v", err)
			continue
		}
		if !wasChanged {
			conf.Log.Info("Base wasn't changed")
			continue
		}

		conf.Log.Info("Parsing new base...")
		err = parseAndStoreGeoBase(conf)
		if err != nil {
			conf.Log.Warnf("Parse error: %v", err)
			continue
		}

		conf.Log.Info("Base was updated")
	}
}

func passiveUpdateRoutine(conf *Config) {
	inf, _ := os.Stat(conf.Dir.baseFile())
	mTime := inf.ModTime()
	clock := time.Tick(conf.UpdInterval)
	for _ = range clock {
		conf.Log.Info("Looking for new base...")

		inf, err := os.Stat(conf.Dir.baseFile())
		if err != nil {
			conf.Log.Warnf("Check error: %v", err)
			continue
		}
		if mTime == inf.ModTime() {
			conf.Log.Info("Base wasn't changed")
			continue
		}

		conf.Log.Info("Parsing new base...")
		err = parseAndStoreGeoBase(conf)
		if err != nil {
			conf.Log.Warnf("Parse error: %v", err)
			continue
		}

		conf.Log.Info("Base was updated")
		mTime = inf.ModTime()
	}
}
