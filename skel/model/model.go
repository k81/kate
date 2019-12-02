package model

import (

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/k81/orm"
	"go.uber.org/zap"

	"__PROJECT_DIR__/config"
)

var logger *zap.Logger

// Init initialize the model setting.
func Init(l *zap.Logger) {
	conf := config.DB

	logger = l
	orm.Debug = true
	orm.SetLogger(logger.With(zap.String("tag", "debug_sql")))
	orm.RegisterDB("default", "mysql", conf.DataSource, conf.MaxIdleConns, conf.MaxOpenConns, conf.ConnMaxLifetime)
}
