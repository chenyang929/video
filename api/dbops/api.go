package dbops

import (
	"database/sql"
	"log"
	"time"
	"video/api/defs"
	"video/api/utils"
)

var (
	userInsert = "INSERT INTO users (login_name, pwd) VALUES (?, ?)"
	pwdSelect  = "SELECT pwd FROM users WHERE login_name = ?"
	userDel    = "DELETE FROM users WHERE login_name = ? and pwd = ?"

	videoInsert = "INSERT INTO video_info (id, author_id, name, display_ctime) VALUES (?, ?, ?, ?)"
	videoSelect = "SELECT author_id, name, display_ctime FROM video_info WHERE id = ?"
	videoDel    = "DELETE FROM video_info WHERE id = ?"

	commentInsert = "INSERT INTO comments (id, video_id, author_id, content) VALUES (?, ?, ?, ?)"
	commentSelect = `SELECT comments.id, comments.content, users.login_name FROM comments 
	INNER JOIN users ON comments.author_id = users.id WHERE comments.video_id = ? and 
	comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)`
)

func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare(userInsert)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare(pwdSelect)
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	defer stmtOut.Close()

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", nil
	}

	return pwd, nil
}

func DelUser(loginName string, pwd string) error {
	stmtOut, err := dbConn.Prepare(userDel)
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	defer stmtOut.Close()
	_, err = stmtOut.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	return nil
}

func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {

	vid := utils.UUIDNew()
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")

	stmtIns, err := dbConn.Prepare(videoInsert)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(videoSelect)
	var aid int
	var name string
	var dct string

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}
	return res, nil
}

func DelVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare(videoDel)
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

func AddNewComments(vid string, aid int, content string) error {
	id := utils.UUIDNew()

	stmtIns, err := dbConn.Prepare(commentInsert)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comments, error) {
	stmtOut, err := dbConn.Prepare(commentSelect)
	var res []*defs.Comments

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}
	defer stmtOut.Close()
	for rows.Next() {
		var id, content, author string
		if err := rows.Scan(&id, &content, &author); err != nil {
			return res, err
		}
		c := &defs.Comments{Id: id, VideoId: vid, Author: author, Content: content}
		res = append(res, c)
	}
	return res, nil
}
