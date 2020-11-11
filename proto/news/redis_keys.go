package news

import "fmt"

// 频道的订阅列表
// hash数据结构，hkey 订阅者ID hval 订阅时间戳
func KeyOfSubListOfChannel(channelName string) string {
	return fmt.Sprintf("news:channel:%s:subscriber", channelName)
}

