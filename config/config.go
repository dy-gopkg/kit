package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"gopkg.in/yaml.v2"
)

type Configs struct {
	Registry struct {
		Addr    string `yaml:"addr"`
	}

	Srv struct {
		SrvName string `yaml:"srvName"`
		SrvId   uint32 `yaml:"srvId"`
		Addr    string `yaml:"addr"`
		Version string `yaml:"version"`
	}

	Log struct {
		Level        string `yaml:"level"`
		FileSize     int    `yaml:"fileSize"`
		FileSizeUnit string `yaml:"fileSizeUnit"`
		JsonFile     bool   `yaml:"jsonFile"`
	}
}



func LoadConfig(filename string, cfgIns *Configs) {

	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		goto FAILED
	}

	err = yaml.Unmarshal(buff, cfgIns)
	if err != nil {
		goto FAILED
	}
	return
FAILED:
	fmt.Printf("LoadConfig failed:%v", err)
	os.Exit(1)
}


