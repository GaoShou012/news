package news

import (
	"google.golang.org/protobuf/proto"
	proto_news "im/proto/news"
	"runtime"
	"sync"
)

type Codec struct {
	donePool sync.Pool
	jobs     chan interface{}
}

func (c *Codec) Run() {
	c.donePool.New = func() interface{} {
		return make(chan error, 1)
	}
	c.jobs = make(chan interface{})

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				jobs := <-c.jobs
				switch jobs.(type) {
				case *JobNewsItemStreamEncode:
					job := jobs.(*JobNewsItemStreamEncode)
					data, err := proto.Marshal(job.input)
					if err != nil {
						job.err <- err
						break
					}
					job.output = data
					job.err <- nil
				case *JobNewsItemStreamDecode:
					job := jobs.(*JobNewsItemStreamDecode)
					item := &proto_news.NewsItem{}
					if err := proto.Unmarshal(job.input, item); err != nil {
						job.err <- err
						break
					}
					job.output = item
					job.err <- nil
				}
			}
		}()
	}
}

func (c *Codec) NewsItemStreamEncode(data *proto_news.NewsItem) ([]byte, error) {
	job := &JobNewsItemStreamEncode{
		input:  data,
		output: nil,
		err:    c.donePool.Get().(chan error),
	}
	c.jobs <- job
	err := <-job.err
	return job.output, err
}

func (c *Codec) NewsItemStreamDecode(data []byte) (*proto_news.NewsItem, error) {
	job := &JobNewsItemStreamDecode{
		input:  data,
		output: nil,
		err:    c.donePool.Get().(chan error),
	}
	c.jobs <- job
	err := <-job.err
	return job.output, err
}

func (c *Codec) 