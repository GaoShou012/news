package news

import (
	"container/list"
	"im/src/frontier"
)

type Client struct {
	Conn      frontier.Conn
	IsCaching bool
	NewsCache *list.List
	//News      *News
	Channels  *Channels
}

//func (c *Client) IsNewMessageId(key string, id string) bool {
//	lastMessageId, ok := c.LastMessageId[key]
//	if !ok {
//		return true
//	}
//	return im.IsNewMessageId(lastMessageId, id)
//}
//func (c *Client) SetLastMessageId(key string, id string) {
//	c.LastMessageId[key] = id
//}

func (c *Client) NewsCachePush(item *Item) {
	c.NewsCache.PushBack(item)
}
func (c *Client) NewsCachePop() *Item {
	ele := c.NewsCache.Front()
	if ele == nil {
		return nil
	}
	c.NewsCache.Remove(ele)
	return ele.Value.(*Item)
}
func (c *Client) Sender(message []byte) {
	c.Conn.Sender(message)
}
