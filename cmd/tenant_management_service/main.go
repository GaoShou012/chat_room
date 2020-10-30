package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"wchatv1/router"
	"wchatv1/services/tenant_management"
)

//go:generate wire ../../services/tenant_management/

var (
	prometheusAddr string
	ginMode        string
)

var loadFlags = micro.Action(func(c *cli.Context) error {
	prometheusAddr = c.String("metrics_address")
	ginMode = c.String("gin_mode")
	return nil
})

func init() {

}

func main() {
	service := micro.NewService(
		micro.Flags(
			&cli.StringFlag{Name: "prometheus_address", Usage: "The prometheus service"},
			&cli.StringFlag{Name: "gin_mode", Usage: "The gin engine mode that is running"},
		),
	)
	service.Init(loadFlags)

	serviceAddr := service.Server().Options().Address
	web := tenant_management.NewHttpService(
		newMysql(),
		newRedis(),
	)
	go router.Run(serviceAddr, web, ginMode)

	fmt.Println("running")
	fmt.Println("Listen:", serviceAddr)

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

func newRedis() *redis.ClusterClient {
	addr := []string{
		"192.168.56.101:9001",
		"192.168.56.101:9002",
		"192.168.56.101:9003",
		"192.168.56.101:9004",
		"192.168.56.101:9005",
		"192.168.56.101:9006",
	}
	password := ""
	minIdleConns := runtime.NumCPU() * 10
	poolSize := runtime.NumCPU() * 20

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              addr,
		MaxRedirects:       0,
		ReadOnly:           false,
		RouteByLatency:     false,
		RouteRandomly:      false,
		ClusterSlots:       nil,
		OnNewNode:          nil,
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           password,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		NewClient:          nil,
		PoolSize:           poolSize,
		MinIdleConns:       minIdleConns,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	})
	return client
}
func newMysql() *gorm.DB {
	return nil
}
