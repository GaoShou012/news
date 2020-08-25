package news

import "im/src/frontier"

type Client struct {
	Conn    frontier.Conn
	SubList []string
}
