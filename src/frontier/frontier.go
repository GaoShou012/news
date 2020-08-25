package frontier

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/golang/glog"
	"net"
	"runtime"
	"sync"
	"time"
	"im/src/ider"
	"im/src/netpoll"
)

type Frontier struct {
	Id int
	Ln net.Listener
	Protocol
	Handler

	Heartbeat        bool
	HeartbeatTimeout int64

	// id pool
	idPool *ider.IdPool

	// poller
	desc   *netpoll.Desc
	poller netpoll.Poller

	// conn event's handler
	// insert,update,delete
	event struct {
		procNum     int
		timestamp   int64
		pool        sync.Pool
		anchor      map[int]*list.Element
		connections []*list.List
		cache       []chan *connEvent
	}

	accept struct {
		procNum int
		cache   chan net.Conn
	}

	epoll struct {
		procNum int
		cache   chan *conn
	}

	message struct {
		procNum int
		pool    sync.Pool
		cache   chan *Message
	}
}

func (f *Frontier) Init() error {
	// 初始化poller
	poller, err := netpoll.New(nil)
	if err != nil {
		return err
	}
	desc := netpoll.Must(netpoll.HandleListener(f.Ln, netpoll.EventRead|netpoll.EventOneShot))
	f.poller = poller
	f.desc = desc

	f.idPool = &ider.IdPool{}
	f.idPool.Init(1000000)

	// 事件处理者
	// 事件会存在异步情况
	f.eventHandler()

	f.accept.procNum = runtime.NumCPU() * 100
	f.accept.cache = make(chan net.Conn, 100000)
	f.onAccept()

	f.epoll.procNum = runtime.NumCPU()
	f.epoll.cache = make(chan *conn, 100000)
	f.onEpoll()

	f.message.procNum = runtime.NumCPU()
	f.message.cache = make(chan *Message, 100000)
	f.onMessage()

	// init driver
	f.Protocol.OnInit(f)

	return nil
}

func (f *Frontier) Start() error {
	return f.poller.Start(f.desc, func(event netpoll.Event) {
		var err error
		defer func() {
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					fmt.Errorf("accept error:%v\n", err)
					time.Sleep(time.Millisecond * 5)
				}
			}
			f.poller.Resume(f.desc)
		}()

		netConn, err := f.Ln.Accept()
		if err != nil {
			return
		}

		f.accept.cache <- netConn
	})
}

// Stop 停止长连接服务
func (f *Frontier) Stop() error {
	if err := f.poller.Stop(f.desc); err != nil {
		return err
	}
	return f.Ln.Close()
}

func (f *Frontier) onAccept() {
	for i := 0; i < f.accept.procNum; i++ {
		go func() {
			for {
				netConn := <-f.accept.cache

				id, err := f.idPool.Get()
				if err != nil {
					netConn.Close()
					glog.Errorln(err)
					continue
				}
				conn := &conn{
					id:             id,
					fid:            f.Id,
					state:          0,
					broken:         false,
					url:            nil,
					conn:           netConn,
					context:        nil,
					connectionTime: time.Now(),
					deadline:       0,
					desc:           nil,
					ConnWriter:     nil,
					ConnReader:     nil,
					ConnCloser:     nil,
				}
				if err := f.Protocol.OnAccept(conn); err != nil {
					f.idPool.Put(id)
					netConn.Close()
					glog.Errorln(err)
					continue
				}

				desc := netpoll.Must(netpoll.HandleRead(netConn))
				conn.desc = desc

				if err := f.Handler.OnOpen(conn); err != nil {
					f.idPool.Put(id)
					netConn.Close()
					glog.Errorln(err)
					continue
				}

				f.epoll.cache <- conn
			}
		}()
	}
}

func (f *Frontier) onEpoll() {
	procNum := runtime.NumCPU()
	for i := 0; i < procNum; i++ {
		go func() {
			for {
				conn := <-f.epoll.cache
				f.eventPush(ConnEventTypeInsert, conn)

				// Here we can read some new message from connection.
				// We can not read it right here in callback, because then we will
				// block the poller's inner loop.
				// We do not want to spawn a new goroutine to read single message.
				// But we want to reuse previously spawned goroutine.
				f.poller.Start(conn.desc, func(event netpoll.Event) {
					go func() {
						data, err := conn.Reader()
						if err != nil {
							f.Handler.OnClose(conn)
							f.poller.Stop(conn.desc)
							f.idPool.Put(conn.id)
							f.eventPush(ConnEventTypeDelete, conn)
						}
						if bytes.Equal(data, []byte("ping")) {
							f.eventPush(ConnEventTypeUpdate, conn)
							return
						}
						message := f.message.pool.Get().(*Message)
						message.conn, message.data = conn, data
						f.message.cache <- message
						f.eventPush(ConnEventTypeUpdate, conn)
					}()
				})
			}
		}()
	}
}

func (f *Frontier) onMessage() {
	f.message.pool.New = func() interface{} {
		return new(Message)
	}
	for i := 0; i < f.message.procNum; i++ {
		go func() {
			for {
				message := <-f.message.cache
				conn, data := message.conn, message.data
				f.Handler.OnMessage(conn, data)
				f.message.pool.Put(message)
			}
		}()
	}
}

func (f *Frontier) eventPush(evt int, conn *conn) {
	if evt == ConnEventTypeInsert {
		conn.deadline = f.event.timestamp + f.HeartbeatTimeout
	} else if evt == ConnEventTypeUpdate {
		if (conn.deadline - f.event.timestamp) < 5 {
			return
		}
		conn.deadline = f.event.timestamp + f.HeartbeatTimeout
	}

	event := f.event.pool.Get().(*connEvent)
	event.Type, event.Conn = evt, conn
	blockNum := conn.id & f.event.procNum
	f.event.cache[blockNum] <- event
}
func (f *Frontier) eventHandler() {
	procNum := runtime.NumCPU()

	go func() {
		ticker := time.Tick(time.Second)
		f.event.timestamp = time.Now().Unix()
		for {
			now := <-ticker
			f.event.timestamp = now.Unix()
		}
	}()
	f.event.pool.New = func() interface{} {
		return new(connEvent)
	}
	f.event.anchor = make(map[int]*list.Element)
	f.event.connections = make([]*list.List, procNum)
	f.event.cache = make([]chan *connEvent, procNum)
	for i := 0; i < procNum; i++ {
		f.event.connections[i] = list.New()
		f.event.cache[i] = make(chan *connEvent, 50000)

		connections := f.event.connections[i]
		events := f.event.cache[i]
		go func(connections *list.List, events chan *connEvent) {
			defer func() { panic("exception error") }()
			tick := time.NewTicker(time.Second)
			for {
				select {
				case now := <-tick.C:
					if f.Heartbeat == false {
						return
					}
					deadline := now.Unix()
					for {
						ele := connections.Front()
						if ele == nil {
							break
						}

						conn := ele.Value.(*conn)
						if conn.deadline > deadline {
							break
						}
						f.Handler.OnClose(conn)
						f.Protocol.OnClose(conn)
						f.poller.Stop(conn.desc)
						f.idPool.Put(conn.id)

						connId := conn.id
						conn.state = connStateWasClosed

						anchor, ok := f.event.anchor[connId]
						if !ok {
							panic(fmt.Sprintf("The conn id is not exists. %d\n", connId))
						}
						delete(f.event.anchor, connId)
						connections.Remove(anchor)
					}
					break
				case event, ok := <-events:
					conn := event.Conn
					connId := conn.id
					switch event.Type {
					case ConnEventTypeInsert:
						if conn.state == connStateIsNothing {
							conn.state = connStateIsWorking
						} else {
							break
						}
						anchor := connections.PushBack(conn)
						f.event.anchor[connId] = anchor
						break
					case ConnEventTypeDelete:
						if conn.state == connStateIsWorking {
							conn.state = connStateWasClosed
						} else {
							break
						}
						anchor, ok := f.event.anchor[connId]

						if !ok {
							panic(fmt.Sprintf("delete the conn id is not exists %d\n", connId))
						}
						delete(f.event.anchor, connId)
						connections.Remove(anchor)
						break
					case ConnEventTypeUpdate:
						if conn.state != connStateIsWorking {
							break
						}
						anchor, ok := f.event.anchor[connId]
						if !ok {
							panic(fmt.Sprintf("update the conn id is not exists %d\n", connId))
						}
						connections.Remove(anchor)
						anchor = connections.PushBack(conn)
						f.event.anchor[connId] = anchor
						break
					}

					if ok {
						f.event.pool.Put(event)
					}
					break
				}
			}
		}(connections, events)
	}
}
