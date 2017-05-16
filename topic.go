package Milena

import (
	"fmt"
	"github.com/JodeZer/Milena/log"
	"gopkg.in/Shopify/sarama.v1"
	"time"
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
	stopC             stopChan
}

func newTopicWorker(consumer sarama.Consumer, metaDB metaStorageEngine, conf *topicWorkerConfig) *topicWorker {
	t := &topicWorker{c: conf, mEngine: metaDB}
	t.topickey = t.c.ClusterNameBelong + "." + t.c.TopicName
	t.topicFileName = t.c.DataDir + "/" + t.c.TopicName + ".log"
	t.consumer = consumer
	t.partitionSettings = make(map[int32]*partionSetting, len(conf.Partions))
	for _, p := range conf.Partions {
		t.partitionSettings[p.Partition] = &p
	}
	t.tEngine = newTopicSimpleStorageEngine(&topicStorageConfig{t.topicFileName})
	t.stopC = make(stopChan, 1)
	return t
}

func (t *topicWorker) Run() {
	ps, err := t.consumer.Partitions(t.c.TopicName)
	if err != nil {
		log.Errorf("%s", err)
		log.Degbugf("%s", err)
		return
	}
	t.partionKeys = make(map[int32]string, len(ps))

	for _, p := range ps {
		dbOff := t.mEngine.GetOffset(genPartitionkey(t.topickey, p))
		pconf, ok := t.partitionSettings[p]
		if !ok {
			pconf = &partionSetting{Start: 0}
		}
		if pconf.Start < dbOff {
			pconf.Start = dbOff
		}
		pc, err := t.consumer.ConsumePartition(t.c.TopicName, p, pconf.Start)
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
	for i := 0; i < len(t.pConsumer)+1; i++ {
		t.stopC <- stopSig{}
	}
	//wait all exit
	time.Sleep(1 * time.Second)
	t.mEngine.Close()
	t.tEngine.Close()
}

func (t *topicWorker) readLoop() {
	msgChan := make(chan *sarama.ConsumerMessage, 1000)
	func() {
		for i, pc := range t.pConsumer {
			go func(spc sarama.PartitionConsumer, index int) {
				for {

					select {
					case msg := <-spc.Messages():
						msgChan <- msg //TODO timeout
					case <-t.stopC:
						spc.Close()
						log.Degbugf("%s~%d rcv stop sig and stop", t.topickey, index)
						return
					case err := <-spc.Errors():
						log.Errorf("%s %s", t.topickey, err)
					}
				}
			}(pc, i)
		}
	}()
	for {
		select {
		case msg := <-msgChan:
			log.Degbugf(genLineMsg(msg)) //TODO storage
			if err := t.tEngine.Append(msg); err != nil {
				log.Errorf("append fail %s", err)
				continue
			}
			t.mEngine.UpdateOffset(t.partionKeys[msg.Partition], msg.Offset+1)
		case <-t.stopC:
			log.Degbugf("%s rcv stop sig and stop write", t.c.TopicName)
			return
		}
	}
}

func genPartitionkey(topicKey string, p int32) string {
	return topicKey + "~" + fmt.Sprintf("%d", p)
}
