// +build wireinject

package news

import (
	"github.com/google/wire"
	"google.golang.org/protobuf/proto"
	proto_im "im/proto/im"
)

func Request(message proto.Message) []byte {
	wire.Build(
		wire.Value("news"),
		proto_im.NewMessage,
		proto_im.Encode,
	)
	return nil
}

func Response(message proto.Message) []byte {
	wire.Build(
		wire.Value("news"),
		proto_im.NewMessage,
		proto_im.Encode,
	)
	return nil
}
