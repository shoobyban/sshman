package backend

import (
	"fmt"
	"log"
	"sync"
)

const (
	InfoPrefix  = "\033[1;34mINF\033[0m "
	ErrorPrefix = "\033[1;31mERR\033[0m "
)

type logEntry struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ILogMode string

var ILogModeClosed ILogMode = ""
var ILogModeStdout ILogMode = "stdout"
var ILogModeWeb ILogMode = "web"

type LogWorker struct {
	Source chan interface{}
	Quit   chan struct{}
}

type ILog struct {
	mode    ILogMode
	workers map[LogWorker]bool
	l       *sync.Mutex
}

func NewLog(web bool) *ILog {
	ilog := ILog{
		l:       &sync.Mutex{},
		workers: map[LogWorker]bool{},
	}
	if !web {
		ilog.mode = ILogModeStdout
	}
	return &ilog
}

func (l *ILog) Open(worker LogWorker) {
	log.Println("Log opened from frontend")
	l.l.Lock()
	defer l.l.Unlock()
	l.workers[worker] = true
}

func (l *ILog) Close(worker LogWorker) {
	log.Println("Frontend log closed")
	l.l.Lock()
	defer l.l.Unlock()
	delete(l.workers, worker)
}

func (l *ILog) AddEntry(e logEntry) {
	l.l.Lock()
	defer l.l.Unlock()
	if e.Type == "error" {
		log.Println(ErrorPrefix + e.Message)
	} else {
		log.Println(InfoPrefix + e.Message)
	}
	for worker := range l.workers {
		worker.Source <- e
	}
}

func (l *ILog) Errorf(msg string, args ...interface{}) {
	l.AddEntry(logEntry{Type: "error", Message: fmt.Sprintf(msg, args...)})
}
func (l *ILog) Infof(msg string, args ...interface{}) {
	l.AddEntry(logEntry{Type: "info", Message: fmt.Sprintf(msg, args...)})
}
