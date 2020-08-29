package news

import (
	"container/list"
	proto_news "im/proto/news"
	"im/src/im"
)

type SubClients map[int]*Client

type Task struct {
	id  int
	ele *list.Element
}

type PullNewsTask struct {
	Client  *Client
	SubList []string
}

type PublisherTask struct {
	Client *Client
	Items  []*Item
}
type PubNewsToClient struct {
	Client *Client
	Items  []*proto_news.NewsItem
}

type News map[string]*Item

func (n News) IsSub(key string) bool {
	_, ok := n[key]
	return ok
}
func (n News) Sub(key string) {
	n[key] = nil
}
func (n News) IsNew(newItem *Item) bool {
	oldItem, ok := n[newItem.Key]
	if !ok {
		return true
	}
	return im.IsNewMessageId(oldItem.Id, newItem.Id)
}
func (n News) GetItem(key string) (*Item, bool) {
	item, ok := n[key]
	return item, ok
}
func (n News) SetItem(item *Item) {
	n[item.Key] = item
}
