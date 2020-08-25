package frontier

type Handler interface {
	OnOpen(conn Conn) error
	OnMessage(conn Conn, message []byte)
	OnClose(conn Conn)
}
