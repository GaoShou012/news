package queue

import (
	"container/list"
	"sync"
)

type Queue struct {
	mutex sync.Mutex
	li    list.List
}

func (q *Queue) Push(v interface{}) *list.Element {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.li.PushBack(v)
}
func (q *Queue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	ele := q.li.Front()
	q.li.Remove(ele)
	return ele.Value
}
func (q *Queue) GetByAnchor(ele *list.Element) {

}
