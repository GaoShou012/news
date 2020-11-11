package logger

import (
	"fmt"
	"testing"
)

func TestLogger_Println(t *testing.T) {
	fmt.Println("测试日志打印")

	logger := Logger{Lv: LvAll}
	logger.Println(LvInfo, "i am info")
	logger.Println(LvWarning, "i am warning")
	logger.Println(LvError, "i am error")
}

func TestLogger_Compare(t *testing.T) {
	fmt.Println("测试日志匹配打印")

	logger := Logger{Lv: LvWarning}
	logger.Compare(LvError, "i am error")
	logger.Compare(LvWarning, "i am warning")
	logger.Compare(LvInfo, "i am info")
}
