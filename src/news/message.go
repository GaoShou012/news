package news

import (
	"im/src/im"
)

type Message struct {
	Type string
	Key  string
	Val  string
}

type Item struct {
	Id  string
	Key string
	Val string
}

func ResponseNewsItems(items []Item) *im.Message {
	msg := &im.Message{
		Head: im.Head{BusinessType: "", BusinessApi: ""},
		Body: items,
	}
	return msg
}


