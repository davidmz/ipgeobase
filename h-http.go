package main

import (
	"encoding/json"
	"net/http"
)

type HttpHandler struct {
	*Config
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	base := h.VBase.Load().(*GeoBase)
	result := base.Find(ip)
	w.Header().Set("Server", ServerName)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}
