package Milena

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type plock struct {
	FileName string
}

func (p *plock) Lock() error {

	if _, err := os.Stat(p.FileName); !os.IsNotExist(err) {
		return errors.New("Process has started")
	}

	if file, err := os.Create(p.FileName); err != nil {
		return fmt.Errorf("create lock failed %s", err)
	} else {
		pid := os.Getpid()
		file.Write([]byte(strconv.Itoa(pid)))
		file.Close()
	}
	return nil
}

func (p *plock) Unlock() {
	os.Remove(p.FileName)
}
