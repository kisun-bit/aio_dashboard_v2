package integration

import "github.com/kisun-bit/aio_dashboard/internal/depends/postgresql"

type IntegratorUnion struct {
	DB postgresql.GetCloser
}
