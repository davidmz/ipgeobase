package main

import (
	"os"
	"path"

	"github.com/davidmz/gotools/mustbe"
)

type DataDir struct{ dirName string }

func (d *DataDir) baseFile() string    { return path.Join(d.dirName, "base.zip") }
func (d *DataDir) eTagFile() string    { return path.Join(d.dirName, "etag.txt") }
func (d *DataDir) tmpBaseFile() string { return d.baseFile() + ".tmp" }
func (d *DataDir) tmpETagFile() string { return d.eTagFile() + ".tmp" }

func (d *DataDir) cleanTemps() (eerr error) {
	defer mustbe.Catched(&eerr)
	mustbe.OK(os.Remove(d.tmpBaseFile()))
	mustbe.OK(os.Remove(d.tmpETagFile()))
	return
}

func (d *DataDir) moveTemps() (eerr error) {
	defer mustbe.Catched(&eerr)
	mustbe.OK(os.Remove(d.baseFile()))
	mustbe.OK(os.Remove(d.eTagFile()))
	mustbe.OK(os.Rename(d.tmpBaseFile(), d.baseFile()))
	mustbe.OK(os.Rename(d.tmpETagFile(), d.eTagFile()))
	return
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
