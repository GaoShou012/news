package frontier

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io/ioutil"
	"net"
	"time"
)

var _ Protocol = &ProtocolWs{}

type ProtocolWs struct {
	frontier             *Frontier
	WriterBufferSize     int
	ReaderBufferSize     int
	WriterMessageTimeout time.Duration
	ReaderMessageTimeout time.Duration

	writer ConnWriter
	reader ConnReader
	closer ConnCloser
}

func (p *ProtocolWs) OnInit(f *Frontier) {
	p.frontier = f

	p.writer = func(conn net.Conn, message []byte) error {
		if err := conn.SetWriteDeadline(time.Now().Add(p.WriterMessageTimeout)); err != nil {
			return err
		}
		w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
		if _, err := w.Write(message); err != nil {
			return err
		}
		return w.Flush()
	}

	p.reader = func(conn net.Conn) ([]byte, error) {
		h, r, err := wsutil.NextReader(conn, ws.StateServerSide)
		if err != nil {
			return nil, err
		}

		if h.OpCode.IsControl() {
			if h.OpCode == ws.OpPing {
				err := conn.SetWriteDeadline(time.Now().Add(p.ReaderMessageTimeout))
				if err != nil {
					return nil, err
				}
				w := wsutil.NewControlWriter(conn, ws.StateServerSide, ws.OpPong)
				if _, e := w.Write(nil); e != nil {
					err = e
				}
				if e := w.Flush(); e != nil {
					err = e
				}
				return nil, nil
			}
			return nil, wsutil.ControlFrameHandler(conn, ws.StateServerSide)(h, r)
		}

		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	p.closer = func(conn net.Conn) error {
		w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpClose)
		if _, err := w.Write([]byte("close")); err != nil {
			return err
		}
		return w.Flush()
	}
}

func (p *ProtocolWs) OnAccept(conn Conn) error {
	if p.frontier.Debug {
		_, err := ws.Upgrader{
			ReadBufferSize:  p.ReaderBufferSize,
			WriteBufferSize: p.WriterBufferSize,
			Protocol:        nil,
			ProtocolCustom:  nil,
			Extension:       nil,
			ExtensionCustom: nil,
			Header:          nil,
			OnRequest: func(uri []byte) error {
				conn.SetUrl(uri)
				fmt.Println("conn uri:", string(uri))
				return nil
			},
			OnHost: func(host []byte) error {
				fmt.Println("conn host", string(host))
				return nil
			},
			OnHeader: func(key, value []byte) error {
				fmt.Println("conn header", key, string(value))
				return nil
			},
			OnBeforeUpgrade: nil,
		}.Upgrade(conn.GetConn())
		if err != nil {
			return nil
		}
	} else {
		_, err := ws.Upgrader{
			ReadBufferSize:  p.ReaderBufferSize,
			WriteBufferSize: p.WriterBufferSize,
			Protocol:        nil,
			ProtocolCustom:  nil,
			Extension:       nil,
			ExtensionCustom: nil,
			Header:          nil,
			OnRequest: func(uri []byte) error {
				conn.SetUrl(uri)
				return nil
			},
			OnHost:          nil,
			OnHeader:        nil,
			OnBeforeUpgrade: nil,
		}.Upgrade(conn.GetConn())
		if err != nil {
			return err
		}
	}

	conn.Init(p.writer, p.reader, p.closer)
	return nil
}
func (p *ProtocolWs) OnClose(conn Conn) error {
	w := wsutil.NewWriter(conn.GetConn(), ws.StateServerSide, ws.OpClose)
	if _, err := w.Write([]byte("close")); err != nil {
		return err
	}
	return w.Flush()
}
