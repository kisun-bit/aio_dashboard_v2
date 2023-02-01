package systemd

import (
	"errors"
	"github.com/kisun-bit/aio_dashboard/internal/depends"
	"github.com/kisun-bit/aio_dashboard/internal/middleware"
	"github.com/kisun-bit/aio_dashboard/pkg/core"
	"go.uber.org/zap"
)

type BackendServer struct {
	Depend depends.Dependency
	Middle middleware.Middleware
	HTTP   core.HTTPMixin
}

func NewBackendServer(globalLogger, cronLogger *zap.SugaredLogger) (*BackendServer, error) {
	if globalLogger == nil || cronLogger == nil {
		return nil, errors.New("logger required")
	}

}
