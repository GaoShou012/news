package news

import "im/src/frontier"

type EventUploadSubscribe struct {
	Client        *Client
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
