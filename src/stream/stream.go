package stream


type Stream interface {
	LastMessageId() string
	Push() error
	GetByMessageId() (interface{},error)
}
