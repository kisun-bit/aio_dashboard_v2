package main

import (
	"os"

	"github.com/kisun-bit/aio_dashboard/configs"
	"github.com/kisun-bit/aio_dashboard/internal/systemd"
	"github.com/kisun-bit/aio_dashboard/pkg/logger"
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
	gLogger := logger.NewLoggerWithDefaultOptions(configs.Settings.Base.Name, configs.Settings.Base.GlobalLogPath)
	cLogger := logger.NewLoggerWithDefaultOptions(configs.Settings.Base.Name, configs.Settings.Base.CronLoggerPath)

	defer func() {
		_ = gLogger.Sync()
		_ = cLogger.Sync()
	}()

	srv, eNewSrv := systemd.NewDashboardSrv(gLogger, cLogger)
	if eNewSrv != nil {
		gLogger.Fatalf("init dashboard service: %v", eNewSrv)
	}

	inst := systemd.ParseInstFromArgs(os.Args...)
	eResp := systemd.ResponseInst(srv, inst)
	if eResp != nil {
		gLogger.Fatalf("control dashboard service [%v]: %v", inst, eResp)
	}
}
