package database

import (
	dbconfig "StructForge/backend/common/data/database/config"

	mysqlAdapter "StructForge/backend/common/data/database/adapters/mysql"
	postgresAdapter "StructForge/backend/common/data/database/adapters/postgres"
	sqliteAdapter "StructForge/backend/common/data/database/adapters/sqlite"
)

// NewDatabase 根据配置创建数据库实例
// 注意：此函数从 init.go 中移出，以避免导入循环
func NewDatabase(config *dbconfig.Config) (Database, error) {
	switch config.AdapterType {
	case dbconfig.AdapterPostgreSQL:
		if config.PostgreSQL == nil {
			config.PostgreSQL = DefaultConfig().PostgreSQL
		}
		return postgresAdapter.NewAdapter(config)
	case dbconfig.AdapterMySQL:
		if config.MySQL == nil {
			config.MySQL = DefaultConfig().MySQL
		}
		return mysqlAdapter.NewAdapter(config)
	case dbconfig.AdapterSQLite:
		if config.SQLite == nil {
			config.SQLite = DefaultConfig().SQLite
		}
		return sqliteAdapter.NewAdapter(config)
	default:
		return nil, ErrNotSupported
	}
}
