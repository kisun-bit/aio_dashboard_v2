package depends

import (
	"github.com/kisun-bit/aio_dashboard/internal/depends/postgresql"
	"github.com/kisun-bit/aio_dashboard/internal/depends/redis"
)

type Dependency struct {
	DB    postgresql.GetCloser
	Cache redis.Operator
}
