package systemd

import (
	"github.com/kardianos/service"
	"github.com/kisun-bit/aio_dashboard/configs"
	"go.uber.org/zap"
	"strings"
	"time"
)

type SrvCtlInstruction string

const (
	Install   SrvCtlInstruction = "install"
	Uninstall SrvCtlInstruction = "uninstall"
	Start     SrvCtlInstruction = "start"
	Stop      SrvCtlInstruction = "stop"
)

type Systemctl struct {
	globalLogger,
	cronLogger *zap.SugaredLogger
	srv *BackendServer
}

func NewDashboardSrv(globalLogger, cronLogger *zap.SugaredLogger) (service.Service, error) {
	srvConfig := &service.Config{
		Name:         configs.Settings.Base.Name,
		DisplayName:  configs.Settings.Base.DisplayName,
		Description:  configs.Settings.Base.Description,
		Dependencies: strings.Split(configs.Settings.Base.SrvDepends, ","),
	}

	ctl := new(Systemctl)
	ctl.globalLogger = globalLogger
	ctl.cronLogger = cronLogger

	return service.New(ctl, srvConfig)
}

func ParseInstFromArgs(args ...string) SrvCtlInstruction {
	if len(args) <= 1 {
		return Install
	}
	return SrvCtlInstruction(args[1])
}

// ResponseInst 响应服务控制指令
// 当inst为非法指令(未被定义)时，以Run指令执行
func ResponseInst(srv service.Service, inst SrvCtlInstruction) error {
	switch inst {
	case Install:
		return srv.Install()
	case Uninstall:
		return srv.Uninstall()
	case Start:
		return srv.Start()
	case Stop:
		return srv.Stop()
	default:
		return srv.Run()
	}
}

func (control *Systemctl) Start(service.Service) error {
	if service.Interactive() {
		control.globalLogger.Info("running in terminal")
	} else {
		control.globalLogger.Info("running under service manager")
	}

	go control.run()
	return nil
}

func (control *Systemctl) run() {
	// TODO 执行服务
}

func (control *Systemctl) Stop(service.Service) (err error) {
	time.Sleep(2 * time.Second)
	return err
}
