package taskrunner

import (
	"errors"
	"log"
	"os"
	"sync"
	"video/scheduler/dbops"
)

func delVideo(vid string) error {
	err := os.Remove(VideoDir + vid)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Del video error: %v", err)
		return err
	}
	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDelRecord(3)
	if err != nil {
		log.Printf("Video clear dispatch error: %v", err)
	}
	if len(res) == 0 {
		return errors.New("All tasks finished")
	}

	for _, id := range res {
		dc <- id
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
forLoop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) {
				if err := delVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := dbops.DelVideoDelRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)

		default:
			break forLoop
		}
	}
	var err error
	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})
	return err
}
