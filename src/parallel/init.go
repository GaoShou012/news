package parallel

import "runtime"

var Codec *Parallel

func init() {
	Codec = &Parallel{}
	Codec.Init(runtime.NumCPU())
}
