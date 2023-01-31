package log

import (
	"fmt"
	"github.com/kisun-bit/aio_dashboard/configs"
	"github.com/kisun-bit/aio_dashboard/pkg/env"
	"github.com/kisun-bit/aio_dashboard/pkg/logger"
	"github.com/kisun-bit/aio_dashboard/pkg/timeutil"
	"go.uber.org/zap"
)

func NewCronLogger() *zap.SugaredLogger {
	return logger.NewLogger(
		logger.WithDisableConsole(),
		logger.WithField("env", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileRotation(logger.GenerateDefaultLBJWriter(configs.ProjectCronLogFile)))
}
