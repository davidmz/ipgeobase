package main

import (
	"encoding/json"
	"strconv"

	"github.com/davidmz/ipgeobase/mmc"
)

const MemcacheVersion = "1.0.0"

type MemcacheHandler struct {
	*Config
}

func (m *MemcacheHandler) ServeMemcache(req *mmc.Request, resp *mmc.Response) error {
	m.Log.Debugf("Memcache command %q", req)

	switch req.Command {

	case "get", "gets": // retrieve commands
		if len(req.Args) == 0 {
			resp.ClientError("key required")
			break
		}
		if req.Command == "get" {
			req.Args = req.Args[:1]
		}
		base := m.VBase.Load().(*GeoBase)
		for _, ip := range req.Args {
			data, _ := json.Marshal(base.Find(ip))
			if err := resp.Value(ip, data); err != nil {
				return err
			}
		}
		resp.Status("END")

	case "set", "add", "replace", "append", "prepend": // store commands
		if len(req.Args) < 4 {
			resp.ClientError("invalid command format")
			break
		}
		bodyLen, err := strconv.Atoi(req.Args[3])
		if err != nil || bodyLen < 0 {
			resp.ClientError("invalid data length")
			break
		}
		req.ReadBody(bodyLen)
		resp.Status("STORED")

	case "version":
		resp.Status("VERSION " + MemcacheVersion)

	case "quit":
		return mmc.CloseConnError

	default:
		resp.UnknownCommandError()
	}

	return nil
}
