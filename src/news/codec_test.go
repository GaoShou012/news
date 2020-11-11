package news

import (
	"fmt"
	"github.com/golang/glog"
	proto_news "im/proto/news"
	"testing"
)

func TestCodec_NewsItemStreamEncode(t *testing.T) {
	fmt.Println("测试Codec NewsItemStreamEncode")
	codec := &Codec{}
	codec.Run()
	data, err := codec.NewsItemStreamEncode(&proto_news.NewsItem{Key: "123"})
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println("encode data:", data)
}
