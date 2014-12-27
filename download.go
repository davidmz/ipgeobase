package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/davidmz/gotools/mustbe"
)

func downloadGeoBase(conf *Config) (wasChanged bool, eerr error) {
	dir := conf.Dir

	defer func() {
		if wasChanged && eerr == nil {
			dir.moveTemps()
		}
		dir.cleanTemps()
	}()
	defer mustbe.Catched(&eerr, &wasChanged)

	req, err := http.NewRequest("GET", conf.GeoBaseUrl, nil)
	mustbe.OK(err)
	if b, err := ioutil.ReadFile(dir.eTagFile()); err == nil {
		req.Header.Set("If-None-Match", string(b))
	}

	resp, err := http.DefaultClient.Do(req)
	mustbe.OK(err)
	defer resp.Body.Close()

	mustbe.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified).
		Else(fmt.Errorf("Http error: %s", resp.Status))

	if resp.StatusCode == http.StatusNotModified {
		// База не поменялась
		return
	}

	wasChanged = true

	mustbe.OK(ioutil.WriteFile(dir.tmpETagFile(), []byte(resp.Header.Get("ETag")), 0666))

	file, err := os.Create(dir.tmpBaseFile())
	mustbe.OK(err)
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	mustbe.OK(err)

	return
}
