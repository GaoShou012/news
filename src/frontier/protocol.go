package frontier

type Protocol interface {
	OnInit(f *Frontier)
	OnAccept(conn Conn) error
	OnClose(conn Conn) error
}
