package logger

import (
	"fmt"
	"runtime"
)

var Lv int

func println(lv int, a ...interface{}) {
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

func Error(a ...interface{}) {
	println(LvError, a...)
}
func Warning(a ...interface{}) {
	println(LvWarning, a...)
}
func Info(a ...interface{}) {
	println(LvInfo, a...)
}

func Compare(lv int, a ...interface{}) {
	if lv > Lv {
		return
	}
	println(lv,a...)
}
