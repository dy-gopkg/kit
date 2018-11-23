package util

import (
	"fmt"
	"github.com/dy-gopkg/kit/log"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-config/source/file"
	"github.com/sirupsen/logrus"
	"os"
)

type AddrAndPath struct {
	Addr string `json:"Addr"`
	Path string `json:"Path"`
}

type BaseConfig struct {
	ServiceConfig AddrAndPath `json:"ServiceConfig"`
	BusinessConfig AddrAndPath `json:"BusinessConfig"`
}

type ServiceConfig struct {
	Registry struct {
		Address string
	}

	Service struct {
		Name string
		ListenAddr string
		BrokerAddr string
		Version string
	}

	Log struct {
		Path string
		Level string
		FileSize int32
	}
}

var (
	BaseConf BaseConfig
	ServiceConf ServiceConfig
)

func LoadConfig() {
	// 加载最基础的配置
	err := config.Load(file.NewSource(file.WithPath("config.json")))
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	err = config.Scan(&BaseConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loadServiceConfig()
}

func loadServiceConfig() {
	// default load from local config file
	err := config.Load(file.NewSource(file.WithPath("service_config.json")))
	if err != nil {
		err = config.Load(consul.NewSource(consul.WithAddress(BaseConf.ServiceConfig.Addr)))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	err = config.Get(BaseConf.ServiceConfig.Path).Scan(&ServiceConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initLogger()

	go watchServiceConfig()
}

func initLogger() {
	// log
	setLoggerLevel()

	logrus.SetOutput(log.NewLogFile(
		log.FilePath(ServiceConf.Log.Path),
		log.FileSize(config.Get("log","fileSize").Int(10),
			config.Get("log","fileUnit").String("M")),
		log.FileTime(true)))

}

func setLoggerLevel() {
	level, _ := logrus.ParseLevel(ServiceConf.Log.Level)
	logrus.SetLevel(level)
}

func watchServiceConfig() {
	w, err := config.Watch(BaseConf.ServiceConfig.Path, "log", "level")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		v, err := w.Next()
		if err != nil {
			fmt.Println(err)
			continue
		}

		ServiceConf.Log.Level = v.String("info")

		setLoggerLevel()
	}

}
