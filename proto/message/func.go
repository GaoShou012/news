package message

import (
	"strconv"
	"strings"
)

func IdToInt(messageId string) (timestamp int, num int, err error) {
	arr := strings.Split(messageId, "-")
	timestamp, err = strconv.Atoi(arr[0])
	if err != nil {
		return
	}
	num, err = strconv.Atoi(arr[1])
	if err != nil {
		return
	}

	return
}

func IsNewMessageId(oldMsgId string, newMsgId string) bool {
	t0, n0, err := IdToInt(oldMsgId)
	if err != nil {
		return false
	}

	t1, n1, err := IdToInt(newMsgId)
	if err != nil {
		return false
	}

	if t1 < t0 {
		return false
	}

	if t1 == t0 && n1 <= n0 {
		return false
	}

	return true
}
