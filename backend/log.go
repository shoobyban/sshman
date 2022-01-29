package backend

import (
	"fmt"
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

type ILog struct {
	l        *sync.RWMutex
	mode     ILogMode
	entries  []logEntry
	notEmpty *sync.Cond
}

func NewLog(web bool) *ILog {
	ilog := ILog{
		l: &sync.RWMutex{},
	}
	if !web {
		ilog.mode = ILogModeStdout
	} else {
		ilog.notEmpty = sync.NewCond(ilog.l.RLocker())
	}
	return &ilog
}

func (l *ILog) isEmpty() bool {
	return len(l.entries) == 0
}

func (l *ILog) GetMode() ILogMode {
	if l == nil { // testing mostly where we created config storage manually
		l = NewLog(false)
	}
	l.l.RLock()
	defer l.l.RUnlock()
	return l.mode
}

func (l *ILog) OpenWeb() {
	l.l.Lock()
	defer l.l.Unlock()
	l.mode = ILogModeWeb
}

func (l *ILog) AddEntry(e logEntry) {
	m := l.GetMode()
	switch m {
	case ILogModeWeb:
		l.l.Lock()
		defer l.l.Unlock()
		l.entries = append(l.entries, e)
		l.notEmpty.Signal()
	case ILogModeStdout:
		if e.Type == "error" {
			fmt.Println(ErrorPrefix + e.Message)
		} else {
			fmt.Println(InfoPrefix + e.Message)
		}
	}
}

func (l *ILog) Errorf(msg string, args ...interface{}) {
	if l.GetMode() != ILogModeClosed {
		l.AddEntry(logEntry{Type: "error", Message: fmt.Sprintf(msg, args...)})
	}
}
func (l *ILog) Infof(msg string, args ...interface{}) {
	if l.GetMode() != ILogModeClosed {
		l.AddEntry(logEntry{Type: "info", Message: fmt.Sprintf(msg, args...)})
	}
}

func (l *ILog) Pop() logEntry {
	if l.isEmpty() {
		l.l.RLock()
		l.notEmpty.Wait()
	}
	l.l.Lock()
	defer l.l.Unlock()
	e := l.entries[0]
	l.entries = l.entries[1:]
	return e
}
