package main

import (
	"net/http"
	"video/scheduler/dbops"

	"github.com/julienschmidt/httprouter"
)

func vidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	if len(vid) == 0 {
		sendResponse(w, 400, "video id is empty")
		return
	}

	err := dbops.AddVideoDelRecord(vid)
	if err != nil {
		sendResponse(w, 500, "Internal server error")
		return
	}
	sendResponse(w, 200, "")
}
