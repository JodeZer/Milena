package Milena

import (
	"gopkg.in/Shopify/sarama.v1"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/JodeZer/Milena/log"
	"encoding/binary"
)

//========================================
//storage Engine
type topicStorageEngine interface {
	Append(sarama.ConsumerMessage) error
}

type metaStorageEngine interface {
	UpdateOffset(string, int64) error
	GetOffset(string) int64
	Close()
}
//=======================================
//implement topicStorage Engine
type topicStorageConfig struct {

}

type topicSimpleStorage struct {

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

func (m *metaLevelDBStorageEngine)UpdateOffset(key string, offset int64) error {
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
