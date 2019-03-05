package kit

import (
	"flag"
	"fmt"
	"github.com/dy-gopkg/kit/util"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	GitTag    = "2000.01.01.release"
	BuildTime = "2000-01-01T00:00:00+0800"
)

var (
	DefaultService micro.Service
)

func Init() {
	//显示版本号信息
	version := flag.Bool("v", false, "version")
	flag.Parse()

	if *version {
		fmt.Println("Git Tag: " + GitTag)
		fmt.Println("Build Time: " + BuildTime)
		os.Exit(0)
	}

	util.LoadConfig()
	DefaultService = micro.NewService(
		micro.Name(util.ServiceConf.Name),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(util.ServiceConf.Version),
		micro.Metadata(util.ServiceConf.Metadata),
		micro.Registry(registry.NewRegistry(registry.Addrs(util.RegistryConf.Address))),
	)

	DefaultService.Init()

}

func Run() {
	if err := DefaultService.Run(); err != nil {
		logrus.Fatalf("service run error: %v", err)
	}
}

func Server() server.Server {
	return DefaultService.Server()
}

func ServiceName() string {
	return util.ServiceConf.Name
}

func ServiceListenAddr() string {
	return util.ServiceConf.ListenAddr
}

func ServiceExternalAddr() string {
	return util.ServiceConf.ExternalAddr
}

func ServiceBrokerAddr() string {
	return util.ServiceConf.BrokerAddr
}

func ServiceVersion() string {
	return util.ServiceConf.Version
}

func ServiceMetadata(key, def string) string {
	val, ok := util.ServiceConf.Metadata[key]
	if ok {
		return val
	}
	return def
}

func Client() client.Client {
	return DefaultService.Client()
}
