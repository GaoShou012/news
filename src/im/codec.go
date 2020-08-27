package im

import (
	"encoding/json"
)

func Encode(message *Message) (data []byte, err error) {
	data, err = json.Marshal(message)
	return
}
