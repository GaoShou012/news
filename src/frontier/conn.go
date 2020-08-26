package frontier

import (
	"im/src/netpoll"
	"net"
	"time"
)

const (
	connStateIsNothing = iota
	connStateIsWorking
	connStateWasClosed
)

type ConnWriter func(conn net.Conn, message []byte) error
type ConnReader func(conn net.Conn) ([]byte, error)
type ConnCloser func(conn net.Conn) error
type ConnPong func(conn net.Conn) error

type Conn interface {
	Init(writer ConnWriter, reader ConnReader, close ConnCloser)
	GetId() int
	GetFid() int
	GetState() int
	GetConnectionTime() time.Time
	Broken()
	IsBroken() bool
	IsWorking() bool
	GetConn() net.Conn
	SetContext(ctx interface{})
	GetContext() interface{}
	SetUrl(url []byte)
	GetUrl() []byte
	Writer(message []byte) error
	Reader() ([]byte, error)
	Closer() error
}

var _ Conn = &conn{}

type conn struct {
	acl map[string]bool

	id     int
	fid    int
	state  int
	broken bool
	url    []byte

	conn           net.Conn
	context        interface{}
	connectionTime time.Time

	deadline int64
	desc     *netpoll.Desc

	ConnWriter
	ConnReader
	ConnCloser
}

func (c *conn) Init(writer ConnWriter, reader ConnReader, closer ConnCloser) {
	c.ConnWriter, c.ConnReader, c.ConnCloser = writer, reader, closer
}

func (c *conn) IsWorking() bool {
	if c.state == connStateIsWorking {
		return true
	} else {
		return false
	}
}

func (c *conn) Writer(message []byte) error {
	return c.ConnWriter(c.conn, message)
}
func (c *conn) Reader() ([]byte, error) {
	return c.ConnReader(c.conn)
}
func (c *conn) Closer() error {
	return c.ConnCloser(c.conn)
}

func (c *conn) Acl(key string) bool {
	if c.acl == nil {
		return false
	}
	_, ok := c.acl[key]
	return ok
}
func (c *conn) SetAcl(acl map[string]bool) {
	c.acl = acl
}
func (c *conn) GetAcl() map[string]bool {
	return c.acl
}

func (c *conn) GetId() int {
	return c.id
}
func (c *conn) GetFid() int {
	return c.fid
}

func (c *conn) GetState() int {
	return c.state
}

func (c *conn) GetConnectionTime() time.Time {
	return c.connectionTime
}

func (c *conn) Broken() {
	c.broken = true
}
func (c *conn) IsBroken() bool {
	return c.broken
}

func (c *conn) GetConn() net.Conn {
	return c.conn
}

func (c *conn) SetContext(ctx interface{}) {
	c.context = ctx
}

func (c *conn) GetContext() interface{} {
	return c.context
}

func (c *conn) SetUrl(url []byte) {
	c.url = url
}
func (c *conn) GetUrl() []byte {
	return c.url
}
