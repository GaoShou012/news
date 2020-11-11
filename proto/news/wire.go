// +build wireinject

package news

import (
	"github.com/google/wire"
	proto_im "gitlab.cptprod.net/wchat-backend/im-lib/proto/im"
	"google.golang.org/protobuf/proto"
)

func Pack(message proto.Message) []byte {
	wire.Build(
		wire.Value("news"),
		proto_im.NewMessage,
		proto_im.Encode,
	)
	return nil
}
