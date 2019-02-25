package taskrunner

import (
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("Dispatcher sent: %d", i)
		}
		return nil
	}

	e := func(dc dataChan) error {
	forLoop:
		for {
			select {
			case dt := <-dc:
				log.Printf("Executor received: %v", dt)
			default:
				break forLoop
			}
		}
		return nil //  errors.New("err")则一次收发完就结束了
	}
	runner := NewRunner(30, false, d, e)
	go runner.StartAll() // 没有go，测试永远不会结束
	time.Sleep(3 * time.Second)
}
