package pubsub

import (
	"container/list"
	"sync"
)

type Queue struct {
	mutex  sync.Mutex
	anchor map[int]*list.Element
	li     *list.List
}

func (q *Queue) Init() {
	q.anchor = make(map[int]*list.Element)
	q.li = list.New()
}

func (q *Queue) Front() *list.Element {
	return q.li.Front()
}
func (q *Queue) Back() *list.Element {
	return q.li.Back()
}

func (q *Queue) Push(index int, v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	anchor := q.li.PushBack(v)
	q.anchor[index] = anchor
}
func (q *Queue) Pop() *list.Element {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	ele := q.li.Front()
	if ele == nil {
		return nil
	}
	q.li.Remove(ele)
	return ele
}

func (q *Queue) Remove(index int) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	anchor, ok := q.anchor[index]
	if !ok {
		return ok
	}
	q.li.Remove(anchor)
	return ok
}
