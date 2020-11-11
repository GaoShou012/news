package news

import "github.com/google/wire"

var Provider = wire.NewSet(
	wire.Struct(new(ChannelNews), "*"),
)
