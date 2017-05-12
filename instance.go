package Milena

import (
	"github.com/JodeZer/Milena/log"
	"os"
)

type Instance struct {
	proclock *plock
}

func NewInsatnce(c *Config) (*Instance, error) {
	ins := &Instance{}
	ins.proclock = &plock{c.LockFile}

	return ins, nil
}

func (i *Instance) Start() {
	//create lock file
	if err := i.proclock.Lock(); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}

}

func (i *Instance) Stop() {
	i.proclock.Unlock()
}