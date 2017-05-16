package Milena

import (
	"gopkg.in/Shopify/sarama.v1"
	"fmt"
	"github.com/JodeZer/Milena/log"
	"time"
	"sync/atomic"
	"sync"
)

//=======================================
//cluster Instance

var omitTopic = []string{
	"_",
}

var (
	running int32 = 1
	stop int32 = 2
	connecting int32 = 3
)

type kafkaClusterConfig struct {
	ClusterName  string
	Brokers      []string
	ListenTopics []topicSetting
	DataDir      string
}
type kafkaCluster struct {
	rw       sync.RWMutex
	c        *kafkaClusterConfig
	tWorkers []*topicWorker
	consumer sarama.Consumer
	mEngine  metaStorageEngine
	mdataDir string
	state    int32
}

func newKafkaCluster(conf *kafkaClusterConfig) (*kafkaCluster, error) {
	if len(conf.Brokers) == 0 {
		return nil, fmt.Errorf("%s has no broker", conf.ClusterName)
	}

	k := &kafkaCluster{c:conf}
	k.state = connecting
	k.mdataDir = k.c.DataDir + "/" + "meta"
	//
	var err error
	if k.mEngine, err = newMetaLevelDBStorageEngine(k.mdataDir); err != nil {
		return nil, err
	}

	return k, nil
}

func (k *kafkaCluster) connectloop() {
	k.rw.Lock()
	defer k.rw.Unlock()
	c := sarama.NewConfig()
	c.Net.DialTimeout = 5 * time.Second
	recon:
	for {
		if atomic.LoadInt32(&k.state) == stop {
			return
		}
		consumer, err := sarama.NewConsumer(k.c.Brokers, c)
		if err != nil {
			log.Errorf("cluster=>[%s] connected fail to %v ,err=>", k.c.ClusterName, k.c.Brokers, err)
			time.Sleep(5 * time.Second)
			continue
		}
		atomic.StoreInt32(&k.state, running)
		k.consumer = consumer
		goto conSucc
	}
	conSucc:
	if len(k.c.ListenTopics) == 0 {
		ts, err := k.consumer.Topics()
		if err != nil {
			goto recon
		}
		ts = removeCommonTopics(ts, omitTopic)
		for _, t := range ts {
			k.c.ListenTopics = append(k.c.ListenTopics, topicSetting{Name:t})
		}

	}

	for _, t := range k.c.ListenTopics {
		tw := newTopicWorker(k.consumer, k.mEngine, &topicWorkerConfig{
			TopicName:t.Name,
			Partions:t.Partitions,
			DataDir:k.c.DataDir,
			ClusterNameBelong:k.c.ClusterName,
		})
		k.tWorkers = append(k.tWorkers, tw)
	}

	for _, tw := range k.tWorkers {
		go tw.Run()
	}

}
func (k *kafkaCluster)Run() {
	go k.connectloop()

}

func (k *kafkaCluster) Stop() {
	atomic.CompareAndSwapInt32(&k.state, connecting, stop)
	k.rw.Lock()
	defer k.rw.Unlock()
	for _, tw := range k.tWorkers {
		log.Degbugf("%s call stop", tw.c.TopicName)
		tw.Stop()
	}
	if k.consumer != nil {
		k.consumer.Close()
	}
	if k.mEngine != nil {
		k.mEngine.Close()
	}
}

func removeCommonTopics(ls []string, rs []string) []string {
	res := make([]string, 0, len(ls))
	m := make(map[string]byte, len(ls))
	for _, v := range ls {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = byte(1)
	}

	for _, v := range rs {
		if _, ok := m[v]; ok {
			delete(m, v)
		}
	}

	for k, _ := range m {
		res = append(res, k)
	}
	return res
}
