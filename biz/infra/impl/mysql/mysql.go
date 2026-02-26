package mysql

import (
	"fmt"

	"github.com/xh-polaris/synapse4b/biz/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	dsn := conf.GetConfig().DB.DSN
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("mysql open, dsn: %s, err: %w", dsn, err)
	}
	return db, nil
}
