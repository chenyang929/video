package session

import (
	"sync"
	"time"
	"video/api/dbops"
	"video/api/defs"
	"video/api/utils"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMill() int64 {
	return time.Now().UnixNano() / 1000000 // 毫秒
}

func delExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DelSession(sid)
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}
	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

func GenerateNewSessionId(un string) string {
	id := utils.UUIDNew()
	ct := nowInMill()
	ttl := ct + 30*60*1000 // 30分钟过期

	ss := &defs.SimpleSession{Username: un, TTL: ttl}
	sessionMap.Store(id, ss)
	dbops.InsertSession(id, ttl, un)
	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := nowInMill()
		if ss.(*defs.SimpleSession).TTL < ct {
			delExpiredSession(sid)
			return "", true
		}
		return ss.(*defs.SimpleSession).Username, false
	}
	return "", true
}
