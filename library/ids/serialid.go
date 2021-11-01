package ids

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	idgenEpoch         = int64(1545998880000)        //开始时间截 (2017-01-01)
	workerIdBits       = uint(10)                    //机器id所占的位数
	maxWorkerId        = -1 ^ (-1 << workerIdBits)   //支持的最大机器id数量
	sequenceBits       = uint(12)                    //序列所占的位数
	workerIdShift      = sequenceBits                //机器id左移位数
	timestampLeftShift = sequenceBits + workerIdBits ///时间戳左移位数
	sequenceMask       = -1 ^ (-1 << sequenceBits)   //4095
)

//
type IdWorker struct {
	sequence      int64
	lastTimestamp int64
	workerId      int64
	mutex         *sync.Mutex
	baseBits      *BaseBits //shorten the id
}

//创建idworker对象
func NewIdWorker(workerId int64) (*IdWorker, error) {
	idWorker := &IdWorker{}
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("illegal worker id: %d", workerId))
	}

	idWorker.workerId = workerId
	idWorker.lastTimestamp = -1
	idWorker.sequence = 0
	idWorker.mutex = &sync.Mutex{}
	baseBits, err := NewBaseBits(int8(62)) //默认radix为62
	if err != nil {
		return nil, errors.New("can not initialize 'baseN4go'")
	}
	idWorker.baseBits = baseBits
	return idWorker, nil
}

//must create
func (id *IdWorker) MustNextId() int64 {
	if id, err := id.NextId(); err == nil {
		return id
	}
	return 0
}

// need synchronized
func (id *IdWorker) NextId() (int64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()

	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		return 0, errors.New(fmt.Sprintf("Clock moved backwards.Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp))
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp
	return ((timestamp - idgenEpoch) << timestampLeftShift) | (id.workerId << workerIdShift) | id.sequence, nil
}

//重置用于shorten id
func (id *IdWorker) RabaseShortRadix(radix int8) error {
	baseBits, err := NewBaseBits(radix)
	if err != nil {
		return err
	}
	id.baseBits = baseBits
	return nil
}

//生成10位短码
func (id *IdWorker) ShortId() (string, error) {
	newId, err := id.NextId()
	if err != nil {
		return "", err
	}
	return id.baseBits.Encode(newId)
}

//根据id生成10位短码
func (id *IdWorker) ShortenId(genId int64) (string, error) {
	return id.baseBits.Encode(genId)
}

//返回id是由哪个workerId生成的
func (id *IdWorker) WorkerId(genId int64) int64 {
	workerId := uint(uint(genId<<42) >> 54)
	return int64(workerId)
}

//
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

//毫秒
func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
