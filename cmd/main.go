package main

import (
	"github.com/kisun-bit/aio_dashboard/configs"
	"github.com/kisun-bit/aio_dashboard/internal/log"
	"github.com/kisun-bit/aio_dashboard/internal/systemd"
	"os"
)

// @title swagger 接口文档
// @version 1.0
// @description 数据备份与恢复管理系统接口文档

// @contact.name aio
// @contact.url
// @contact.email kisun668@gmail.com

// @license.name TODO
// @license.url TODO

// @securityDefinitions.apikey  LoginToken
// @in                          header
// @name                        token

// @BasePath /
func main() {
	globalLogger := log.NewGlobalLogger()
	cronLogger := log.NewCronLogger()

	defer func() {
		_ = globalLogger.Sync()
		_ = cronLogger.Sync()
	}()

	srv, err := systemd.NewDashboardSrv(configs.Settings.Basic, globalLogger, cronLogger)
	if err != nil {
		globalLogger.Fatalf("init dashboard service: %v", err)
	}

	inst := systemd.ParseInstFromArgs(os.Args...)
	err = systemd.ResponseInst(srv, inst)
	if err != nil {
		globalLogger.Fatalf("control dashboard service [%v]: %v", inst, err)
	}
}
