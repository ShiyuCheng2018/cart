package main

import (
	"fmt"
	"github.com/ShiyuCheng2018/cart/domain/repository"
	"github.com/ShiyuCheng2018/cart/domain/service"
	"github.com/ShiyuCheng2018/cart/handler"
	cart "github.com/ShiyuCheng2018/cart/proto/cart"
	"github.com/ShiyuCheng2018/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

var QPS = 100

func main() {
	// consul configuration
	consulConfig, err := common.ConsulConfigurator("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	// consul Register
	consulRegistry := consul.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{
				"127.0.0.1:8500",
			}
		})

	// Tracing Analysis
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()

	opentracing.SetGlobalTracer(t)

	// Database configuration
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	// DB connection
	fmt.Print(mysqlInfo.User + ":" + mysqlInfo.Pwd + "@/" + mysqlInfo.Database + "?charset=utf8&parseTime=True&loc=Local")

	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	// DB initialization
	//rp := repository.NewCartRepository(db)
	//rp.InitTable()

	db.SingularTable(true)

	// New Service
	srv := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		// settings for address & port
		micro.Address("0.0.0.0:8087"),
		// add consul as registry center
		micro.Registry(consulRegistry),
		// add tracing Analysis
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper(opentracing.GlobalTracer())),
		// add Queries per second (QPS) [the amount of search traffic an information-retrieval system]
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// initialization for service
	srv.Init()

	cartDataService := service.NewCartDataService(repository.NewCartRepository(db))

	// Register Handler
	err = cart.RegisterCartHandler(srv.Server(), &handler.Cart{CartDataService: cartDataService})
	if err != nil {
		return
	}

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
