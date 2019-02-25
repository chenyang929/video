package dbops

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	VideoSelect = "SELECT video_id FROM video_del_rec LIMIT ?"
	VideoDel    = "DELETE FROM video_del_rec WHERE video_id = ?"
)

func ReadVideoDelRecord(count int) ([]string, error) {
	stmtOut, err := dbConn.Prepare(VideoSelect)
	ids := []string{}
	if err != nil {
		return ids, err
	}
	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDelRecord error: %v", err)
		return ids, err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}

	defer stmtOut.Close()
	return ids, nil
}

func DelVideoDelRecord(vid string) error {
	stmtDel, err := dbConn.Prepare(VideoDel)
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Del VideoDelRecord error: %v", err)
		return err
	}

	defer stmtDel.Close()
	return nil
}
