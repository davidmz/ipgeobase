package main

import (
	"encoding/binary"
	"net"
	"sort"
)

type GeoBase struct {
	Blocks []*GeoBaseBlock
	Cities map[int]*CityInfo
}

type GeoBaseBlock struct {
	Start    uint32 `json:"-"`
	End      uint32 `json:"-"`
	Interval string `json:"inetnum"`
	Country  string `json:"country"`
	CityID   int    `json:"-"`
}

type CityInfo struct {
	Name     string  `json:"city"`
	Region   string  `json:"region"`
	District string  `json:"district"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type SearchResult struct {
	Ip    string `json:"ip"`
	Error string `json:"error,omitempty"`
	*GeoBaseBlock
	*CityInfo
}

func (s SearchResult) ToMap() (mapa map[string]interface{}) {
	mapa = make(map[string]interface{})
	mapa["ip"] = s.Ip

	if s.GeoBaseBlock == nil {
		mapa["error"] = s.Error
		return
	}

	mapa["inetnum"] = s.Interval
	mapa["country"] = s.Country

	if s.CityInfo == nil {
		return
	}

	mapa["city"] = s.Name
	mapa["region"] = s.Region
	mapa["district"] = s.District
	mapa["lat"] = s.Lat
	mapa["lng"] = s.Lng

	return
}

func NewGeoBase() *GeoBase {
	return &GeoBase{
		Blocks: make([]*GeoBaseBlock, 0),
		Cities: make(map[int]*CityInfo),
	}
}

func (g *GeoBase) Find(ipStr string) (result SearchResult) {
	result.Ip = ipStr
	result.Error = "Not found"
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return
	}

	intAddr := binary.BigEndian.Uint32(ip)
	found := sort.Search(len(g.Blocks), func(i int) bool { return g.Blocks[i].End >= intAddr })
	if found != len(g.Blocks) {
		bl := g.Blocks[found]
		if intAddr >= bl.Start {
			result.GeoBaseBlock = bl
			result.CityInfo = g.Cities[bl.CityID]
			result.Error = ""
		}
	}
	return
}
