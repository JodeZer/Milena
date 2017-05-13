package Milena

import (
	"gopkg.in/Shopify/sarama.v1"
	"github.com/JodeZer/Milena/log"
	"fmt"
)

var iF = func(b bool, l interface{}, r interface{}) interface{} {
	if b {
		return l
	}
	return r
}

type topicWorkerConfig struct {
	TopicName         string
	Partions          []partionSetting
	DataDir           string
	ClusterNameBelong string
}
type topicWorker struct {
	c                 *topicWorkerConfig
	mEngine           metaStorageEngine
	topickey          string
	tEngine           topicStorageEngine
	topicFileName     string
	consumer          sarama.Consumer
	pConsumer         []sarama.PartitionConsumer
	partionKeys       map[int32]string
	partitionSettings map[int32]*partionSetting
}

func newTopicWorker(consumer sarama.Consumer, metaDB metaStorageEngine, conf *topicWorkerConfig) *topicWorker {
	t := &topicWorker{c:conf, mEngine:metaDB}
	t.topickey = t.c.ClusterNameBelong + "." + t.c.TopicName
	t.topicFileName = t.c.DataDir + "/" + t.c.TopicName + ".log"
	t.consumer = consumer
	t.partitionSettings = make(map[int32]*partionSetting, len(conf.Partions))
	for _, p := range conf.Partions {
		t.partitionSettings[p.Partition] = &p
	}
	t.tEngine = newTopicSimpleStorageEngine(&topicStorageConfig{t.topicFileName})
	return t
}

func (t *topicWorker) Run() {
	for k, v := range t.partitionSettings {
		dbOff := t.mEngine.GetOffset(genPartitionkey(t.topickey, k))
		if v.Start < dbOff {
			t.partitionSettings[k].Start = dbOff
		}
	}

	ps, err := t.consumer.Partitions(t.c.TopicName)
	if err != nil {
		log.Errorf("%s", err)
		log.Degbugf("%s", err)
		return
	}
	t.partionKeys = make(map[int32]string, len(ps))

	for _, p := range ps {
		pc, err := t.consumer.ConsumePartition(t.c.TopicName, p, t.partitionSettings[p].Start)
		t.partionKeys[p] = t.topickey + "~" + fmt.Sprintf("%d", p)
		if err != nil {
			log.Errorf("%s", err)
			log.Degbugf("%s", err)
			continue
		}
		t.pConsumer = append(t.pConsumer, pc)
	}
	log.Degbugf("%v", t.partionKeys)
	go t.readLoop()

}

func (t *topicWorker) Stop() {

}

func (t *topicWorker)readLoop() {
	msgChan := make(chan *sarama.ConsumerMessage, 1000)
	func() {
		for _, pc := range t.pConsumer {
			go func(spc sarama.PartitionConsumer) {
				for {
					msg := <-spc.Messages()
					msgChan <- msg//TODO timeout
				}
			}(pc)
		}
	}()
	for msg := range msgChan {
		log.Degbugf(genLineMsg(msg))//TODO storage
		if err := t.tEngine.Append(msg); err != nil {
			log.Errorf("append fail %s", err)
			continue
		}
		_ = t.partionKeys[msg.Partition]
		t.mEngine.UpdateOffset(t.partionKeys[msg.Partition], msg.Offset + 1)
	}
}

func genPartitionkey(topicKey string, p int32) string {
	return topicKey + "~" + fmt.Sprintf("%d", p)
}