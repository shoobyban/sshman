package backend

import (
	"fmt"
	"log"
	"sync"
)

const (
	// InfoPrefix is the prefix for info messages
	InfoPrefix = "\033[1;34mINF\033[0m "
	// ErrorPrefix is the prefix for error messages
	ErrorPrefix = "\033[1;31mERR\033[0m "
)

type logEntry struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ILogMode is the mode for the log
type ILogMode string

// ILogModeClosed is the closed mode
var ILogModeClosed ILogMode = ""

// ILogModeStdout is the stdout mode for CLI
var ILogModeStdout ILogMode = "stdout"

// ILogModeWeb is the web mode for the log
var ILogModeWeb ILogMode = "web"

// LogWorker is log transfer agent for web
type LogWorker struct {
	Source chan interface{}
	Quit   chan struct{}
}

// ILog is a logging engine with ability to send logs to multiple targets
type ILog struct {
	mode    ILogMode
	workers map[LogWorker]bool
	l       *sync.Mutex
}

// NewLog creates a new ILog log engine, if web is false it sets StdOut mode
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

// Open adds a web worker to the log handling a keep-alive connection
func (l *ILog) Open(worker LogWorker) {
	log.Println("Log opened from frontend")
	l.l.Lock()
	defer l.l.Unlock()
	l.workers[worker] = true
}

// Close removes a worker from the log
func (l *ILog) Close(worker LogWorker) {
	log.Println("Frontend log closed")
	l.l.Lock()
	defer l.l.Unlock()
	delete(l.workers, worker)
}

func (l *ILog) print(e logEntry) {
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

// Errorf prints an error message to the log
func (l *ILog) Errorf(msg string, args ...interface{}) {
	l.print(logEntry{Type: "error", Message: fmt.Sprintf(msg, args...)})
}

// Infof prints an info message to the log
func (l *ILog) Infof(msg string, args ...interface{}) {
	l.print(logEntry{Type: "info", Message: fmt.Sprintf(msg, args...)})
}
