package cli

import "flag"

var (
	Start int32 = 1
	Signal int32 = 2
)

type Cli struct {
	// command type
	CliType  int32

	// signal
	Sig      struct {
				 //Signal keyword
				 Keyword string
			 }
	// conf file
	ConfFile string
}

func ParseArgs() *Cli {
	c := &Cli{}
	signal := flag.String("s", "", "signal")
	conf := flag.String("f", "", "conf file")
	flag.Parse()
	if *signal != "" {
		c.CliType = Signal
		c.Sig.Keyword = *signal
		return c
	}

	c.CliType = Start
	c.ConfFile = *conf
	return c
}
