package news

import (
	"fmt"
	proto_news "im/proto/news"
	"runtime"
)

type syncRecord struct {
	channels []string
	jobs     chan *proto_news.Subscribe
}

func (s *syncRecord) init() {
	s.jobs = make(chan *proto_news.Subscribe, 100000)
	for i := 0; i < runtime.NumCPU()*10; i++ {
		go func() {
			for {
				job := <-s.jobs
				fmt.Println(job)
			}
		}()
	}
}

func (s *syncRecord) push(subscribe *proto_news.Subscribe) {
	s.jobs <- subscribe
}
