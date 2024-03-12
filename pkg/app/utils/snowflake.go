package utils

import (
	"fmt"
	"time"
)

type Snowflake struct {
	lastTime     int64
	sequence     int64
	workerID     int64
	workerBits   uint
	sequenceBits uint
	timeShift    uint
	workerShift  uint
	maxWorkerID  int64
	maxSequence  int64
	idChan       chan int64
}

// CreateSnowflake 随机生成一个唯一标识
func CreateSnowflake(workerID int64) int64 {
	//workerID := int64(1)
	workerBits := uint(5)
	sequenceBits := uint(9)
	bufferSize := 100
	sf, err := NewSnowflake(workerID, workerBits, sequenceBits, bufferSize)
	if err != nil {
		return 0
	}
	return sf.NextID()
}

// NewSnowflake 雪花算法生成唯一ID
// workerID     用于标识唯一的 worker 或节点，
// workerBits   用于确定 workerID 所占用的位数
// sequenceBits 用于确定序列号所占用的位数
// bufferSize   用于确定生成ID时的缓冲区大小 - 最大并发量
func NewSnowflake(workerID int64, workerBits uint, sequenceBits uint, bufferSize int) (*Snowflake, error) {
	maxWorkerID := int64((1 << workerBits) - 1)
	maxSequence := int64((1 << sequenceBits) - 1)

	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}

	// 创建一个带缓冲区的通道
	idChan := make(chan int64, bufferSize)

	sf := &Snowflake{
		lastTime:     time.Now().UnixNano() / 1000000,
		workerID:     workerID,
		workerBits:   workerBits,
		sequenceBits: sequenceBits,
		timeShift:    workerBits + sequenceBits,
		workerShift:  sequenceBits,
		maxWorkerID:  maxWorkerID,
		maxSequence:  maxSequence,
		idChan:       idChan,
	}

	// 启动一个goroutine来生成ID并发送到通道中
	go sf.generateIDs()

	return sf, nil
}

func (sf *Snowflake) generateIDs() {
	for {
		if sf.idChan != nil {
			curTime := time.Now().UnixNano() / 1000000

			sf.incrementSequence(curTime)

			id := ((curTime << sf.timeShift) | (sf.workerID << sf.workerShift) | sf.sequence)

			select {
			case sf.idChan <- id:
			default:
				// 如果通道已满，丢弃当前ID
			}
		}
	}
}

func (sf *Snowflake) incrementSequence(curTime int64) {
	if curTime == sf.lastTime {
		sf.sequence = (sf.sequence + 1) & sf.maxSequence
		if sf.sequence == 0 {
			// 暂停一毫秒，等待下一个时间窗口
			time.Sleep(1 * time.Millisecond)
			sf.incrementSequence(time.Now().UnixNano() / 1000000)
		}
	} else {
		sf.sequence = 0
		sf.lastTime = curTime
	}
}

func (sf *Snowflake) NextID() int64 {
	return <-sf.idChan
}
