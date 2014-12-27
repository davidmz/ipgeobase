package main

import (
	"os"
	"path"
)

type DataDir struct{ dirName string }

func (d *DataDir) baseFile() string    { return path.Join(d.dirName, "base.zip") }
func (d *DataDir) eTagFile() string    { return path.Join(d.dirName, "etag.txt") }
func (d *DataDir) tmpBaseFile() string { return d.baseFile() + ".tmp" }
func (d *DataDir) tmpETagFile() string { return d.eTagFile() + ".tmp" }

func (d *DataDir) cleanTemps() {
	os.Remove(d.tmpBaseFile())
	os.Remove(d.tmpETagFile())
}

func (d *DataDir) moveTemps() {
	os.Remove(d.baseFile())
	os.Remove(d.eTagFile())
	os.Rename(d.tmpBaseFile(), d.baseFile())
	os.Rename(d.tmpETagFile(), d.eTagFile())
}

func FileExists(filename string) (ok bool, err error) {
	_, err = os.Stat(filename)
	if err == nil {
		ok = true
	}
	if os.IsNotExist(err) {
		ok = false
		err = nil
	}
	return
}
