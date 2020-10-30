package main

import (
	"github.com/golang/glog"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"os"
	"os/signal"
	"syscall"
	"wchatv1/config"
	"wchatv1/router"
	"wchatv1/utils"
)

var (
	prometheusAddr string
	ginMode        string
)

var loadFlags = micro.Action(func(c *cli.Context) error {
	prometheusAddr = c.String("metrics_address")
	ginMode = c.String("gin_mode")
	return nil
})

func main() {
	service := micro.NewService(
		micro.Flags(
			&cli.StringFlag{Name: "prometheus_address", Usage: "The prometheus service"},
			&cli.StringFlag{Name: "gin_mode", Usage: "The gin engine mode that is running"},
		),
	)
	service.Init(loadFlags)

	utils.Micro.Init(service)
	utils.Micro.LoadSource()
	utils.Micro.LoadConfigMust(config.RoomServiceConfig)

	go utils.Prometheus(prometheusAddr)

	serviceAddr := service.Server().Options().Address
	web := &router.Tenant{}
	go router.Run(serviceAddr, web, ginMode)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch s := <-c; s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			glog.Infof("got signal %s; stop server", s)
		case syscall.SIGHUP:
			glog.Infof("got signal %s; go to deamon", s)
			continue
		}
		break
	}
}
