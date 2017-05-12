package Milena

import (
	"os"
	"github.com/pkg/errors"
	"fmt"
	"strconv"
)

type plock struct {
	FileName string
}

func (p *plock)Lock() error {

	if _, err := os.Stat(p.FileName); !os.IsNotExist(err) {
		return errors.New("Process has started")
	}

	if file, err := os.Create(p.FileName); err != nil {
		return errors.New(fmt.Sprintf("create lock failed %s", err))
	} else {
		pid := os.Getpid()
		file.Write([]byte(strconv.Itoa(pid)))
		file.Close()
	}
	return nil
}

func (p *plock)Unlock() {
	os.Remove(p.FileName)
}



