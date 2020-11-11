package news

import (
	"fmt"
	"sync"
	"time"
)

/*
	频道消息拉取
	puller 每秒进行一次拉取任务，一次最大拉取N个频道的消息
	每个新的连接，接入的时候，都会会把频道投放到puller，一秒内，多个客户接入
	可能订阅的频道都是一致的，所以把客户重复订阅的频道去掉，进行一次性拉取
*/
type Puller struct {
	channels []string

	mutex   sync.Mutex
	bufFlag int
	buf     []map[string]bool
}

func (p *Puller) Add(channel string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.buf[p.bufFlag][channel] = true
}

func (p *Puller) Init() {
	p.buf = make([]map[string]bool, 2)
	p.buf[0] = make(map[string]bool)
	p.buf[1] = make(map[string]bool)

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			p.mutex.Lock()
			bufFlag := p.bufFlag
			if p.bufFlag == 0 {
				p.bufFlag = 1
			} else {
				p.bufFlag = 0
			}
			p.mutex.Unlock()

			// TODO 拉取频道消息
			channels := p.buf[bufFlag]
			fmt.Print("需要拉取消息的频道:", channels)
			for channelName, _ := range p.buf[bufFlag] {
				fmt.Println(channelName)
			}
		}
	}()
}
