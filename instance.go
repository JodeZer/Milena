package Milena

import (
	"github.com/JodeZer/Milena/log"
	"os"
	"sync"
)

type stopSig struct{}
type StopChan chan stopSig
type Instance struct {
	rw       sync.RWMutex
	c        *Config
	proclock *plock
	clusters []*kafkaCluster
}

func NewInsatnce(conf *Config) (*Instance, error) {
	ins := &Instance{}
	ins.proclock = &plock{conf.LockFile}
	ins.c = conf
	return ins, nil
}

func (i *Instance) Start() {

	i.rw.Lock()
	defer i.rw.Unlock()
	//create lock file
	if err := i.proclock.Lock(); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
	for _, s := range i.c.Servers {
		k, err := newKafkaCluster(&kafkaClusterConfig{
			ClusterName:s.Name,
			Brokers:s.Brokers,
			DataDir:s.DataDir,
			ListenTopics:s.Topics,
		})
		if err != nil {
			log.Errorf("start failed %s", err)
			os.Exit(1)
		}
		i.clusters = append(i.clusters, k)
	}

	for _, c := range i.clusters {
		c.Run()
	}

}

func (i *Instance) Stop() {
	i.rw.Lock()
	defer i.rw.Unlock()
	i.proclock.Unlock()
	for _, c := range i.clusters {
		c.Stop()
	}
}

