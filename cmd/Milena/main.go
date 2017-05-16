package main

import (
	"fmt"
	"github.com/JodeZer/Milena"
	"github.com/JodeZer/Milena/cmd/cli"
	"github.com/JodeZer/Milena/log"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	_SIG_STOP = "stop"
)

func main() {
	//TODO -f to point config file
	clia := cli.ParseArgs()
	if clia.CliType == cli.Start {
		_Milena_main(clia.ConfFile)
		return
	} else if clia.CliType == cli.Signal {
		_Milena_SigCommands(clia.Sig.Keyword)
		return
	}

}

func _Milena_main(confFile string) {
	c := Milena.NewConfig(confFile)
	log.Degbugf("%+v", c)
	ins, err := Milena.NewInsatnce(c)
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
	// start listen instance
	ins.Start()

	log.Infof("Milena Starts pid:%d", os.Getpid())

	// rcv close signal
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTSTP)
	//signal.Notify(ch, syscall.Signal(0xff))
	s := <-ch
	log.Infof("rev close signal %s", s.String())

	// gracefully close instance
	ins.Stop()

	//log out
	log.Infof("Milena Passed Away, Farewell Kafka")

	// give time to pass away
	time.Sleep(100 * time.Millisecond)
}

func _Milena_SigCommands(sig string) {
	if sig == _SIG_STOP {
		file, err := os.Open("Milena.lock")
		if err != nil {
			log.Errorf("load process fail %s", err)
			return
		}
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			log.Errorf("read lock file fail %s", err)
		}

		cmdStr := fmt.Sprintf("kill -2 %s", bs)
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		if err := cmd.Start(); err != nil {
			log.Errorf("stop fail %s", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			log.Errorf("stop fail %s", err)
			return
		}
	}
}
