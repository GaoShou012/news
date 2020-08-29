package news

import (
	"im/src/frontier"
	"im/src/im"
)

type UploadSubList struct {
	Client  *Client
	SubList []string
}

type EventUploadSubscribe struct {
	Conn          frontier.Conn
	Message       *im.Message
	SubscribeList []string
}

type EventCleanSubscribe struct {
	Client *Client
}

type EventDownloadSubscribe struct {
	Conn frontier.Conn
}

type EventAddSubList struct {
	Client        *Client
	SubscribeList []string
}
type EventDelSubList struct {
	Client        *Client
	SubscribeList []string
}
