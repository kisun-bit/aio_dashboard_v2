package main

import (
	"fmt"
	"github.com/kisun-bit/aio_dashboard/configs"
	"github.com/kisun-bit/aio_dashboard/pkg/env"
	"github.com/kisun-bit/aio_dashboard/pkg/logger"
	"github.com/kisun-bit/aio_dashboard/pkg/timeutil"
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

	// 初始化全局日志记录器
	globalLogger := logger.NewLogger(
		logger.WithDisableConsole(),
		logger.WithField("env", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileRotation(logger.GenerateDefaultLBJWriter(configs.ProjectAccessLogFile)))

	// 初始化后台日志记录器
	backendLogger := logger.NewLogger(
		logger.WithDisableConsole(),
		logger.WithField("env", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileRotation(logger.GenerateDefaultLBJWriter(configs.ProjectCronLogFile)))

	defer func() {
		_ = globalLogger.Sync()
		_ = backendLogger.Sync()
	}()

	// 依赖注入(postgresql、redis、nats)

}
