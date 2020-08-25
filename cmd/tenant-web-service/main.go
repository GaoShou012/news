package main

import (
	"github.com/micro/go-micro/v2/web"
	"time"
	"wchatv1/config"
	"wchatv1/router"
	"wchatv1/utils"
)

func main() {
	service := web.NewService(
		web.Name(config.TenantWebService.ServiceName()),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*10),
		web.Handler(router.New(&router.Tenant{})),
	)
	service.Init()
	
	utils.Micro.InitV2()
	utils.Micro.LoadSource()


	//if err := utils.Micro.LoadConfig(config.MysqlConfig); err != nil {
	//	panic(err)
	//}
	//fmt.Println("load config:Mysql", fmt.Sprintf("%s:%d", config.MysqlConfig.Addr, config.MysqlConfig.Port))
	//utils.InitMysql(
	//	config.MysqlConfig.User,
	//	config.MysqlConfig.Password,
	//	config.MysqlConfig.Addr,
	//	config.MysqlConfig.Port,
	//	config.MysqlConfig.Database,
	//)

	service.Run()
}
