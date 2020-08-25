package frontier

import (
	"sync"
)

/*
	消息发送器
	所有的群组消息，Push to Sender
	由Sender使用epoll模式，调用应用层Conn进行发送
	message必须是[]byte，避免重复进行消息格式转换
*/

var ParallelSender sender

type job struct {
	conn    Conn
	message []byte
}

type sender struct {
	jobs     []chan *job
	pool     sync.Pool
	parallel int
}

func (s *sender) Init(parallel int, cacheSize int) {
	s.parallel = parallel
	s.jobs = make([]chan *job, parallel)
	s.pool.New = func() interface{} {
		return new(job)
	}
	for i := 0; i < parallel; i++ {
		s.jobs[i] = make(chan *job, cacheSize)
		go func(i int) {
			for {
				job := <-s.jobs[i]
				conn, message := job.conn, job.message
				if conn.IsBroken() == false {
					err := conn.Writer(message)
					if err != nil {
						conn.Broken()
					}
				}
				s.pool.Put(job)
			}
		}(i)
	}
}

func (s *sender) Push(conn Conn, message []byte) {
	j := s.pool.Get().(*job)
	j.conn, j.message = conn, message
	s.jobs[j.conn.GetId()%s.parallel] <- j
}
