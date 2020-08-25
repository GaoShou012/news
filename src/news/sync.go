package news

import (
	"container/list"
	"time"
)

type Sync struct {
	clients *list.List
	task    chan *Client
}

func (s *Sync) Push(cli *Client) {
	s.clients.PushBack(cli)
}
func (s *Sync) Pop() *Client {
	ele := s.clients.Front()
	if ele == nil {
		return nil
	}
	s.clients.Remove(ele)
	return ele.Value.(*Client)
}

func (s *Sync) Handler() {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 200)
		for {
			<-ticker.C
			if cli := s.Pop(); cli == nil {
				continue
			} else {
				s.task <- cli
			}
		}
	}()
	for i := 0; i < 100; i++ {
		go func() {
			for{
				cli := <- s.task

			}
		}()
	}
}
