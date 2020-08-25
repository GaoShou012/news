package news

import "wchatv1/src/frontier"

const (
	EventInsert = iota
	EventDelete
)

type Event struct {
	Type int
	Conn frontier.Conn
}
