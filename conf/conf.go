package conf

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"flag"
	"fmt"
)

type Conf struct {
	Host string
	Port int
	Prefix string
	Db struct {
		Host string
		Port string
		User string
		Password string
		Database string
	}
	Enable bool
	Path string
}

var c Conf

func init() {
	configFile := c.readParams()
	c.readConf(configFile)
}

func (c *Conf) readConf(configFile string) {
	if configFile == "" {
		execPath, err := os.Executable()
		if err != nil {
			log.Printf("os.Executable err  #%v", err)
		}
		configFile = filepath.Join(filepath.Dir(execPath), "../etc/config.yaml")
	}
	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Printf("yamlFile.Get err  #%v", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func (c *Conf) readParams() string{
	var configFile string
	flag.StringVar(&configFile,"c", "", "Specify the config file. The default value is ./etc/conf.yaml")
	//flag.StringVar(&(c.Host), "h", "0.0.0.0", "host name")
	//flag.IntVar(&(c.Port), "p", 8080, "service port")
	flag.Parse()
	flag.Usage()
	fmt.Printf("c=%s, Port=%d\n", configFile, c.Port)
	return configFile
}

func GetConf() Conf {
	return c
}