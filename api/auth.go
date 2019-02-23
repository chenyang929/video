package main

import (
	"net/http"
	"video/api/defs"
	"video/api/session"
)

var (
	headerFieldSession = "X-Session-Id"
	headerFielUname    = "X-User-Name"
)

func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(headerFieldSession)
	if len(sid) == 0 {
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	r.Header.Add(headerFielUname, uname)
	return true
}

func validateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(headerFielUname)
	if len(uname) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}
	return true
}
