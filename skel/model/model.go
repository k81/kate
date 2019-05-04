package model

import (
	"context"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/k81/log"
	"github.com/k81/orm"

	"__PROJECT_DIR__/config"
)

var (
	mctx = log.WithContext(context.Background(), "module", "model")
)

// Init initialize the model setting.
func Init() {
	conf := config.DB

	orm.Debug = true
	orm.SetLogger(log.Tag("__debug_sql"))
	orm.RegisterDB("default", "mysql", conf.DataSource, conf.MaxIdleConns, conf.MaxOpenConns, conf.ConnMaxLifetime)
}
