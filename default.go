package kit

import (
	"flag"
	"fmt"
	"github.com/dy-gopkg/kit/util"
	"github.com/micro/go-config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
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
		micro.Name(config.Get("srv","srvName").String("default")),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Get("srv","version").String("0.0.0")),
		micro.Metadata(map[string]string{"ID": strconv.FormatUint(uint64(config.Get("srv","srvId").Int(0)), 10)}),
		micro.Registry(registry.NewRegistry(registry.Addrs(config.Get("registry","addr").String("0.0.0.0:8500")))))

	DefaultService.Init()


}


func Run(){
	if err := DefaultService.Run(); err != nil {
		logrus.Fatalf("service run error: %v", err)
	}
}
