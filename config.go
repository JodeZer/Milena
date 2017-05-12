package Milena

import (
	"path/filepath"
	"os"
	"github.com/JodeZer/Milena/log"
	"github.com/jinzhu/configor"
)

func init() {
	var err error
	curDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		os.Exit(1)
	}

}

const defaultConfigName = "Milena.yml"
const defaultLockFileName = "Milena.lock"

var curDir string

type Config struct {
	//config file path
	FileName string

	//data files dir
	DataDir  string

	//log level
	LogLevel string

	//lockDir
	LockFile string `yaml:"-"`

	//listen servers
	Servers  []struct {
		//kafka cluster name
		Name        string

		// storage consumer metadata format: {clusterName}/{topic}.meta
		MetaDataDir string`yaml:"-"`

		// storage consumer consumed data format:{}
		DataDir     string`yaml:"-"`

		// Server address
		Brokers     []string


		// listen all topics
		ListenAll   bool

		// listen topic
		Topics      []*topicSetting
	}
}

type topicSetting struct {
	//topic name
	Name  string

	// this time start offset. configed or from metadata
	Start int
}

func NewConfig(fileName string) *Config {
	c := &Config{}
	if fileName == "" {
		fileName = curDir + "/" + defaultConfigName
	}
	c.FileName = fileName

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Errorf("%s  not exists", fileName)
		os.Exit(1)
	}
	log.Infof("find conf %s", fileName)

	c.parse()

	c.valid()

	c.setDefault()

	return c
}

// parse config yml file
func (c *Config) parse() {
	configor.Load(c, c.FileName)
}

// valid if the file is legal else exist
func (c *Config) valid() {
	// valid server conf
	m := make(map[string]string, len(c.Servers))
	for _, server := range c.Servers {
		if server.Name == "" {
			log.Errorf("empty cluser name is not allowed")
			os.Exit(1)
		}

		if len(server.Brokers) == 0 {
			log.Errorf("none broker for %s", server.Name)
			os.Exit(1)
		}

		if server.ListenAll == false && len(server.Topics) == 0 {
			log.Errorf("none topic for %s", server.Name)
			os.Exit(1)
		}

		if _, ok := m[server.Name]; ok {
			log.Errorf("duplicate cluser name is not allowed")
			os.Exit(1)
		}
		m[server.Name] = "1"
	}
	m = nil // helpGC
}

func (c *Config) setDefault() {
	if c.DataDir == "" {
		c.DataDir = curDir + "/" + "data"
	}

	c.LockFile = curDir + "/" + "Milena.lock"
}