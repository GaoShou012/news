package news

import (
	"container/list"
	"im/src/frontier"
)

type ConnList struct {
	anchor      map[int]*list.Element
	connections *list.List
}

func (p *ConnList) Init() {
	p.anchor = make(map[int]*list.Element)
	p.connections = list.New()
}
func (p *ConnList) Push(conn frontier.Conn) {
	anchor := p.connections.PushBack(conn)
	p.anchor[conn.GetId()] = anchor
}
func (p *ConnList) Remove(conn frontier.Conn) {
	anchor, ok := p.anchor[conn.GetId()]
	if !ok {
		return
	}
	p.connections.Remove(anchor)
}
func (p *ConnList) GetConnections() *list.List {
	return p.connections
}
