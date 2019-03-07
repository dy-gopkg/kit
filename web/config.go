package web


import (
	"fmt"
	"github.com/dy-gopkg/kit/log"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-config/source/file"
	"github.com/sirupsen/logrus"
	"os"
)

type (
	RegistryConfig struct {
		Address string
	}

	ServiceConfig struct {
		Name         string
		Version      string
		Metadata     map[string]string
	}

	LogConfig struct {
		Path     string
		Level    string
		FileSize int
	}
)

var (
	DefaultRegistryConf RegistryConfig
	DefaultServiceConf ServiceConfig //服务配置
	DefaultLogConf LogConfig
)

func LoadConfig() {
	// 加载最基础的配置
	err := config.Load(file.NewSource(file.WithPath("config.yaml")))
	if err != nil {
		// load from consul
		fmt.Println("load config from consul")
		addr := os.Getenv("K8S_SERVER_CONFIG_ADDR")
		path := os.Getenv("K8S_SERVER_CONFIG_PATH")
		err = config.Load(consul.NewSource(
			consul.WithAddress(addr),
			consul.WithPrefix(path)))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	// 加载注册中心地址
	err = config.Get("registry").Scan(&DefaultRegistryConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 加载服务配置
	err = config.Get("service").Scan(&DefaultServiceConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 加载并初始化日志配置
	loadLogConfAndInitLogger()
}

func loadLogConfAndInitLogger() {
	err := config.Get("log").Scan(&DefaultLogConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	level, _ := logrus.ParseLevel(DefaultLogConf.Level)
	logrus.SetLevel(level)

	logrus.SetOutput(log.NewLogFile(
		log.FilePath(DefaultLogConf.Path),
		log.FileSize(DefaultLogConf.FileSize),
		log.FileTime(true)))

	go watchLogConf()
}

func watchLogConf() {
	w, err := config.Watch("log")
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

		err = v.Scan(&DefaultLogConf)
		if err != nil {
			continue
		}

		level, _ := logrus.ParseLevel(DefaultLogConf.Level)
		logrus.SetLevel(level)

		logrus.SetOutput(log.NewLogFile(
			log.FilePath(DefaultLogConf.Path),
			log.FileSize(DefaultLogConf.FileSize),
			log.FileTime(true)))
	}

}


