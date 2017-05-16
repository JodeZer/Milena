package cli

import "flag"

var (
	// Start Cmd Type
	Start  int32 = 1

	// Signal Cmd Type
	Signal int32 = 2
)

// Cli
type Cli struct {
	// command type
	CliType int32

	// signal
	Sig struct {
		//Signal keyword
		Keyword string
	}
	// conf file
	ConfFile string
}

// ParseArgs
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
