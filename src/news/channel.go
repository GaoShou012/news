package news

import "container/list"

/*
	频道
	Clients 记录频道下所有订阅的客户
	Developer: 高手
*/
type Channel struct {
	Clients list.List
}

/*
	频道添加一个客户端
*/
func (c *Channel) Add(cli *Client) *list.Element {
	return c.Clients.PushBack(cli)
}

/*
	频道移除一个客户端
*/
func (c *Channel) Del(ele *list.Element) {
	c.Clients.Remove(ele)
}

/*
	频道的客户数量
*/
func (c *Channel) Len() int {
	return c.Clients.Len()
}
