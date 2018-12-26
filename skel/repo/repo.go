package repo

import (
	"context"

	"github.com/k81/kate/log"
	"github.com/k81/kate/orm"
)

var (
	mctx = log.SetContext(context.Background(), "module", "repo")
)

func Init() {
	var (
		confBasic = GetBasicConfig()
		confPools = GetPoolsConfig()

		driverName      = confBasic.DriverName
		dataSource      = confBasic.DataSource
		maxIdle         = confPools.MaxIdleConns
		maxOpen         = confPools.MaxOpenConns
		connMaxLifetime = confPools.ConnMaxLifetime
	)

	orm.Debug = confBasic.DebugSql
	orm.RegisterDB("default", driverName, dataSource, maxIdle, maxOpen, connMaxLifetime)

	//将__DB_NAME__换成对应的数据库名称, 将__TYPE1__ 等换成对应的数据类型的实例.
	//orm.RegisterModel(__DB_NAME__, __TYPE1__, __TYPE2__, ...)
}
