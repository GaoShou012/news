package news

import (
	"container/list"
	"github.com/GaoShou012/frontier"
	proto_news "im/proto/news"
)

type handler struct {
	clients map[int]*Client

	onNews      chan *proto_news.NewsItem
	onLeave     chan frontier.Conn
	onSubscribe chan *EventOnSubscribe

	channels map[string]*list.List
	anchors  map[int]map[string]*list.Element // clientId.channelName.ele
}

func (n *handler) addSubscribe(cli *Client, channels []string) {
	clientId := cli.Conn.GetId()

	for _, channelName := range channels {
		// If the channel is not exists then to create the channel
		channel, ok := n.channels[channelName]
		if !ok {
			channel = list.New()
			n.channels[channelName] = channel
		}

		// Add the client to the client list of the channel
		anchor := channel.PushBack(cli)

		// Save the anchor to the anchor list of the client
		anchorsOfClient, ok := n.anchors[clientId]
		if !ok {
			anchorsOfClient = make(map[string]*list.Element)
			n.anchors[clientId] = anchorsOfClient
		}
		anchorsOfClient[channelName] = anchor
	}
}
func (n *handler) delSubscribe(clientId int) {
	anchorsOfClient, ok := n.anchors[clientId]
	if !ok {
		return
	}

	for channelName, anchor := range anchorsOfClient {
		n.channels[channelName].Remove(anchor)
	}
	delete(n.anchors, clientId)
}

func (n *handler) OnInit() {
	chanSize := 100000
	n.clients = make(map[int]*Client)
	n.onNews = make(chan *proto_news.NewsItem, chanSize)
	n.onSubscribe = make(chan *EventOnSubscribe, chanSize)
	n.onLeave = make(chan frontier.Conn, chanSize)
	n.channels = make(map[string]*list.List)
	n.anchors = make(map[int]map[string]*list.Element)
	go func() {
		for {
			select {
			case event := <-n.onSubscribe:
				clientId, conn, message := event.Conn.GetId(), event.Conn, event.Message

				cli, ok := n.clients[event.Conn.GetId()]
				if !ok {
					cli = &Client{
						Conn:      conn,
						IsCaching: true,
						News:      make(map[string]*proto_news.NewsItem),
					}

					// map[string]*list.Element
					// That is anchors of channel of clients
					n.anchors[clientId] = make(map[string]*list.Element)
					n.clients[clientId] = cli
				}
				n.delSubscribe(clientId)
				n.addSubscribe(cli, event.Message.Channels)

				Agent.onSubscribe <- message.Channels
				break
			case conn := <-n.onLeave:
				clientId := conn.GetId()
				n.delSubscribe(clientId)
				delete(n.clients, clientId)
				break
			case message := <-n.onNews:
				data := proto_news.Response(message)
				for _, cli := range n.clients {
					if cli.IsNewMessageId(message) == false {
						continue
					}
					cli.SetMessageId(message)
					cli.Conn.Sender(data)
				}
				break
			}
		}
	}()
}

func (n *handler) OnSubscribe(conn frontier.Conn, subscribe *proto_news.Subscribe) {
	event := &EventOnSubscribe{
		Conn:    conn,
		Message: subscribe,
	}
	n.onSubscribe <- event
}
func (n *handler) OnLeave(conn frontier.Conn) {
	n.onLeave <- conn
}
func (n *handler) OnNews(item *proto_news.NewsItem) {
	n.onNews <- item
}
