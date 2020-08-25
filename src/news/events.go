package news

import "im/src/frontier"

type EventUploadSubscribe struct {
	Conn          frontier.Conn
	SubscribeList []string
}

type EventCleanSubscribe struct {
	Conn frontier.Conn
}
