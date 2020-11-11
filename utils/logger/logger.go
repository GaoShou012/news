package logger

import (
	"fmt"
	"runtime"
)

const (
	LvClose = iota
	LvError
	LvWarning
	LvInfo
	LvAll
)

type Logger struct {
	Lv int64
}

func (l *Logger) println(lv int64, a ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		switch lv {
		case LvError:
			fmt.Printf("Error:%v %v\n", file, line)
			fmt.Println(a...)
			break
		case LvWarning:
			fmt.Printf("Warning:%v %v\n", file, line)
			fmt.Println(a...)
			break
		case LvInfo:
			fmt.Printf("Info:%v %v\n", file, line)
			fmt.Println(a...)
			break
		}
	}
}

func (l *Logger) Println(lv int64, a ...interface{}) {
	l.println(lv, a...)
}

func (l *Logger) Compare(lv int64, a ...interface{}) {
	if lv > l.Lv {
		return
	}
	l.println(lv, a...)
}

