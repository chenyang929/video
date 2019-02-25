package taskrunner

const (
	ReadyToDispatch = "d"
	ReadToExecute   = "e"
	CLOSE           = "c"
	VideoDir        = "D:/videos/"
)

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error
