package ider

import (
	"container/list"
	"fmt"
	"sync"
)

type IdPool struct {
	mutex sync.Mutex
	li    *list.List
	state map[int]bool
}

func (p *IdPool) Init(size int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.li = list.New()
	p.state = make(map[int]bool)
	for i := 0; i < size; i++ {
		p.li.PushBack(i)
		p.state[i] = true
	}
}

func (p *IdPool) Get() (int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ele := p.li.Front()
	if ele == nil {
		return 0,fmt.Errorf("no id")
	}
	id,ok := ele.Value.(int)
	if !ok {
		return 0,fmt.Errorf("id declare error")
	}
	p.state[id] = false
	p.li.Remove(ele)
	return id,nil
}

func (p *IdPool) Put(id int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	state,ok := p.state[id]
	if !ok {
		return fmt.Errorf("id is error")
	}
	if state {
		return fmt.Errorf("id is exists in pool already")
	}
	p.state[id] = true
	p.li.PushFront(id)
	return nil
}
