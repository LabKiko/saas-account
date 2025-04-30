package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	// 开始时间戳 (2023-01-01 00:00:00 UTC)
	twepoch = int64(1672531200000)

	// 每部分占用的位数
	workerIDBits     = uint(5)  // 机器ID所占的位数
	datacenterIDBits = uint(5)  // 数据中心ID所占的位数
	sequenceBits     = uint(12) // 序列号所占的位数

	// 最大值
	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)     // 最大机器ID
	maxDatacenterID = int64(-1) ^ (int64(-1) << datacenterIDBits) // 最大数据中心ID
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)     // 最大序列号

	// 位移
	workerIDShift      = sequenceBits                      // 机器ID左移位数
	datacenterIDShift  = sequenceBits + workerIDBits       // 数据中心ID左移位数
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits // 时间戳左移位数
)

// Snowflake 结构体
type Snowflake struct {
	mutex         sync.Mutex // 互斥锁
	lastTimestamp int64      // 上次生成ID的时间戳
	workerID      int64      // 机器ID
	datacenterID  int64      // 数据中心ID
	sequence      int64      // 序列号
}

// NewSnowflake 创建一个新的Snowflake实例
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID超出范围")
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("datacenter ID超出范围")
	}
	return &Snowflake{
		lastTimestamp: -1,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}, nil
}

// NextID 生成下一个ID
func (s *Snowflake) NextID() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixNano() / 1000000 // 当前时间戳，精确到毫秒

	// 如果当前时间小于上次生成ID的时间戳，说明系统时钟回退过，应当抛出异常
	if timestamp < s.lastTimestamp {
		return 0, errors.New("时钟回退，拒绝生成ID")
	}

	// 如果是同一时间生成的，则进行序列号递增
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		// 序列号已经达到最大值，等待下一毫秒
		if s.sequence == 0 {
			// 阻塞到下一个毫秒，获得新的时间戳
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 时间戳改变，序列号重置
		s.sequence = 0
	}

	// 更新上次生成ID的时间戳
	s.lastTimestamp = timestamp

	// 移位并通过或运算拼到一起组成64位的ID
	id := ((timestamp - twepoch) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id, nil
}
