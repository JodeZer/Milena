package main

import (
	"github.com/JodeZer/Milena"
	"os"
	"os/signal"
	"github.com/JodeZer/Milena/log"
	"time"
	"syscall"
)

func main() {
	//TODO -f to point config file
	c := Milena.NewConfig("")
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

	signal.Notify(ch, syscall.SIGINT,syscall.SIGTSTP)
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
