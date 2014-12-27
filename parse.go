package main

import (
	"archive/zip"
	"bufio"
	"strconv"
	"strings"

	"github.com/davidmz/go-charset"
	"github.com/davidmz/gotools/mustbe"
)

func parseAndStoreGeoBase(conf *Config) (eerr error) {
	defer mustbe.Catched(&eerr)

	zipFile, err := zip.OpenReader(conf.Dir.baseFile())
	mustbe.OK(err)
	defer zipFile.Close()

	base := NewGeoBase()

	for _, f := range zipFile.File {
		if f.Name == "cidr_optim.txt" {
			// читаем блоки
			rc, err := f.Open()
			mustbe.OK(err)
			func() {
				defer rc.Close()

				scanner := bufio.NewScanner(rc)
				for scanner.Scan() {
					block := &GeoBaseBlock{}
					parts := strings.Split(strings.TrimSpace(scanner.Text()), "\t")

					addr, err := strconv.ParseUint(parts[0], 10, 64)
					mustbe.OK(err)
					block.Start = uint32(addr)

					addr, err = strconv.ParseUint(parts[1], 10, 64)
					mustbe.OK(err)
					block.End = uint32(addr)

					block.Interval = parts[2]

					block.Country = parts[3][:2]
					block.CityID = 0
					if parts[4] != "-" {
						block.CityID, err = strconv.Atoi(parts[4])
						mustbe.OK(err)
					}
					base.Blocks = append(base.Blocks, block)
				}
			}()

		} else if f.Name == "cities.txt" {
			// читаем города
			rc, err := f.Open()
			mustbe.OK(err)
			func() {
				defer rc.Close()

				scanner := bufio.NewScanner(rc)
				for scanner.Scan() {
					line := charset.CP1251.Decode(scanner.Bytes())
					parts := strings.Split(strings.TrimSpace(line), "\t")

					id, err := strconv.Atoi(parts[0])
					mustbe.OK(err)

					cityInfo := &CityInfo{
						Name:     parts[1],
						Region:   parts[2],
						District: parts[3],
					}

					fl, err := strconv.ParseFloat(parts[4], 64)
					mustbe.OK(err)
					cityInfo.Lat = fl

					fl, err = strconv.ParseFloat(parts[5], 64)
					mustbe.OK(err)
					cityInfo.Lng = fl

					base.Cities[id] = cityInfo
				}
			}()

		}
	}

	conf.VBase.Store(base)

	return
}
