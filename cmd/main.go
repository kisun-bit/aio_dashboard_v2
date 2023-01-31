package main

import (
	"github.com/kisun-bit/aio_dashboard/internal/log"
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
	globalLogger := log.NewGlobalLogger()

	// 初始化后台日志记录器
	cronLogger := log.NewCronLogger()

	defer func() {
		_ = globalLogger.Sync()
		_ = cronLogger.Sync()
	}()

	// 依赖注入(postgresql、redis、nats)

}
