package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
	"video/api/defs"
)

var (
	sessionInsert = "INSERT INTO sessions (session_id, TTL, login_name) VALUES (?, ?, ?)"
	sessionSelect = "SELECT TTL, login_name FROM sessions WHERE session_id = ?"
	sessionAll    = "SELECT session_id, TTL, login_name FROM sessions"
	sessionDel    = "DELETE FROM sessions WHERE session_id = ?"
)

func InsertSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare(sessionInsert)
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare(sessionSelect)
	if err != nil {
		return nil, err
	}
	var ttl string
	var uname string
	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.Username = uname
	} else {
		return nil, err
	}

	defer stmtOut.Close()
	return ss, nil
}

func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare(sessionAll)
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	if rows.Next() {
		var id string
		var ttlstr string
		var loginName string
		if err := rows.Scan(&id, &ttlstr, &loginName); err != nil {
			log.Printf("retrieve session error: %s", err)
			return nil, err
		}

		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err == nil {
			ss := &defs.SimpleSession{Username: loginName, TTL: ttl}
			m.Store(id, ss)
			log.Printf("session id: %s, ttl: %d", id, ss.TTL)
		}
	}

	return m, nil
}

func DelSession(sid string) error {
	stmtOut, err := dbConn.Prepare(sessionDel)
	if err != nil {
		log.Printf("%s", err)
	}

	if _, err := stmtOut.Exec(sid); err != nil {
		return nil
	}
	return nil
}
