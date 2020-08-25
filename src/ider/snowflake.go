package ider

import (
	"sync"
	"time"
)

/*
	41bits  毫秒级时间戳 0x1FFFFFFFFFF
	10bits  机器ID
	12bits	每毫秒ID序列号(0-4095)
*/

var SnowFlake *snowFlake

func InitSnowFlake(mid int) {
	SnowFlake = &snowFlake{machineId: int64(mid) << 12}
}

type snowFlake struct {
	mutex         sync.Mutex
	machineId     int64
	sn            int64
	lastTimeStamp int64
}

func (s *snowFlake) Get() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentTimeStamp := time.Now().UnixNano() / 1000000
	if currentTimeStamp == s.lastTimeStamp {
		if s.sn >= 4095 {
			time.Sleep(time.Millisecond)
			currentTimeStamp = time.Now().UnixNano() / 1000000
			s.sn = 0
		}
	}
	s.sn++
	s.lastTimeStamp = currentTimeStamp

	timestamp := currentTimeStamp & 0x1FFFFFFFFFF
	timestamp <<= 22
	id := timestamp | s.machineId | s.sn
	return id
}
