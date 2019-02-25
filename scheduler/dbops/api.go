package dbops

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var VideoInsert = "INSERT INTO video_del_rec (video_id) VALUES (?)"

func AddVideoDelRecord(vid string) error {
	stmtIns, err := dbConn.Prepare(VideoInsert)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDelRecord error: %v", err)
		return err
	}
	defer stmtIns.Close()
	return nil
}
