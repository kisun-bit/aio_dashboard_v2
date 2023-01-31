package systemd

import "github.com/kardianos/service"

// AIODashboardService aio_dashboard以systemd后台运行
type AIODashboardService struct {
	Backend *BackendIntegration // aio_dashboard服务集成
	exit    chan struct{}
}

func (aio *AIODashboardService) Start(s service.Service) (err error) {
	if service.Interactive() {
		// 运行在控制台
	} else {
		// 运行在systemd
	}

	aio.exit = make(chan struct{})

	go aio.run()

	return nil
}

func (aio *AIODashboardService) run() {
	// 执行服务
}

func (aio *AIODashboardService) Stop(s service.Service) (err error) {
	close(aio.exit)
	return err
}
