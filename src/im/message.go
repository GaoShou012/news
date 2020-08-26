package im

import (
	"strconv"
	"strings"
)

type Head struct {
	Id           uint64 `json:"id"`
	BusinessType string `json:"businessType"`
	BusinessApi  string `json:"businessApi"`
}
type Message struct {
	Head `json:"head"`
	Body interface{} `json:"body"`
}

func MessageToInt(id string) (timestamp int, num int, err error) {
	arr := strings.Split(id, "-")
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

func IsNewMessageId(oldId string, newId string) bool {
	t0, n0, err := MessageToInt(oldId)
	if err != nil {
		return false
	}

	t1, n1, err := MessageToInt(newId)
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
