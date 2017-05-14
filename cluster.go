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
	k.state = running
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
	c.Net.DialTimeout = 5 *time.Second
	recon:
	for atomic.LoadInt32(&k.state) == running {
		consumer, err := sarama.NewConsumer(k.c.Brokers, c)
		if err != nil {
			log.Errorf("cluster=>[%s] connected fail to %v ,err=>", k.c.ClusterName, k.c.Brokers, err)
			time.Sleep(5 * time.Second)
			continue
		}
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
	atomic.CompareAndSwapInt32(&k.state, running, stop)
	k.rw.Lock()
	defer k.rw.Unlock()
	for _, tw := range k.tWorkers {
		log.Degbugf("%s call stop", tw.c.TopicName)
		tw.Stop()
	}
	k.consumer.Close()
	k.mEngine.Close()
}

//TODO TODO TODO this is a really dummy ass !!!!!
func removeCommonTopics(ls []string, rs []string) []string {
	f := func(ls []string, r string) []string {
		for i, l := range ls {
			if l == r {
				res := ls[:i]
				res = append(res, ls[i + 1:]...)
				return res
			}
		}
		return ls
	}

	for _, r := range rs {
		ls = f(ls, r)
	}
	return ls
}