package postgresql

import (
	"gorm.io/gorm"
)

type GetCloser interface {
	GetDBForRead() *gorm.DB
	GetDBForWrite() *gorm.DB
	DBRClose() error
	DBWClose() error
}
