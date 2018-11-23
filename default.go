package kit

import (
	"flag"
	"fmt"
	"github.com/dy-gopkg/kit/util"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
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
	ServiceConfigAddr string
	ServiceConfigPath string
	BusinessConfigAddr string
	BusinessConfigPath string
)



func Init(){
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
		micro.Name(util.ServiceConf.Service.Name),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(util.ServiceConf.Service.Version),
		micro.Metadata(util.ServiceConf.Service.Metadata),
		micro.Registry(registry.NewRegistry(registry.Addrs(util.ServiceConf.Registry.Address))))

	DefaultService.Init()


}


func Run(){
	if err := DefaultService.Run(); err != nil {
		logrus.Fatalf("service run error: %v", err)
	}
}
