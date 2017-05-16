package Milena

import (
	"encoding/binary"
	"fmt"
	"github.com/JodeZer/Milena/log"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/Shopify/sarama.v1"
	"os"
	"sync"
)

//========================================
//storage Engine
type topicStorageEngine interface {
	Append(*sarama.ConsumerMessage) error
	Close()
}

type metaStorageEngine interface {
	UpdateOffset(string, int64) error
	GetOffset(string) int64
	Close()
}

//=====================================
//implement metadata storage engine
func newMetaLevelDBStorageEngine(dir string) (metaStorageEngine, error) {
	e := &metaLevelDBStorageEngine{}
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	e.db = db
	return e, nil
}

type metaLevelDBStorageEngine struct {
	db *leveldb.DB
}

func (m *metaLevelDBStorageEngine) UpdateOffset(key string, offset int64) error {
	var b []byte = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(offset))
	m.db.Put([]byte(key), b, nil)
	return nil
}

func (m *metaLevelDBStorageEngine) GetOffset(key string) int64 {
	val, err := m.db.Get([]byte(key), nil)
	if err != nil {
		log.Warnf("fail to load offset for %s err=> %s,default to 0", key, err)
		return int64(0)
	}
	return int64(binary.BigEndian.Uint64(val))
}

func (m *metaLevelDBStorageEngine) Close() {
	m.db.Close()
}

//=======================================
//implement topicStorage Engine
func newTopicSimpleStorageEngine(conf *topicStorageConfig) topicStorageEngine {
	e := &topicSimpleStorageEngine{}
	if f, err := os.OpenFile(conf.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		log.Errorf("%s", err)
	} else {
		e.file = f
	}
	return e
}

type topicStorageConfig struct {
	FileName string
}

type topicSimpleStorageEngine struct {
	file *os.File
	frw  sync.RWMutex
	//buffer *bytes.Buffer
}

func (e *topicSimpleStorageEngine) Append(msg *sarama.ConsumerMessage) error {
	if e.file == nil {
		return errors.New("file open failed")
	}
	e.frw.RLock()
	defer e.frw.RUnlock()
	e.file.WriteString(genLineMsg(msg))
	e.file.Sync()
	return nil
}

func (e *topicSimpleStorageEngine) Close() {
	e.frw.Lock()
	defer e.frw.Unlock()
	e.file.Close()
}

func genLineMsg(msg *sarama.ConsumerMessage) string {
	return fmt.Sprintf("ts=>[%s] p:%d o:%d => %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Partition, msg.Offset, msg.Value)
}
