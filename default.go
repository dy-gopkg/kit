package kit

import (
	"flag"
	"fmt"
	"github.com/dy-gopkg/kit/config"
	"github.com/dy-gopkg/kit/log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)


var (
	GitTag    = "2000.01.01.release"
	BuildTime = "2000-01-01T00:00:00+0800"
)

var (
	DefaultConfig *config.Configs
	DefaultService micro.Service
)



func Init(){
	//显示版本号信息　
	version := flag.Bool("v", false, "version")
	flag.Parse()

	if *version {
		fmt.Println("Git Tag: " + GitTag)
		fmt.Println("Build Time: " + BuildTime)
		return
	}
	DefaultConfig = &config.Configs{}
	config.LoadConfig("config.yaml",DefaultConfig)

	// log
	level, _ := logrus.ParseLevel(DefaultConfig.Log.Level)
	logrus.SetLevel(level)

	logrus.SetOutput(log.NewLogFile(
		log.FilePath("log"),
		log.FileSize(DefaultConfig.Log.FileSize, DefaultConfig.Log.FileSizeUnit),
		log.FileTime(true)))

	DefaultService = micro.NewService(
		micro.Name(DefaultConfig.Srv.SrvName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(DefaultConfig.Srv.Version),
		micro.Metadata(map[string]string{"ID": strconv.FormatUint(uint64(DefaultConfig.Srv.SrvId), 10)}),
		micro.Registry(registry.NewRegistry(registry.Addrs(DefaultConfig.Registry.Addr))))

	DefaultService.Init()


}


func Run(){
	if err := DefaultService.Run(); err != nil {
		logrus.Fatalf("service run error: %v", err)
	}
}
